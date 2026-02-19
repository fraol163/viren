package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fraol163/viren/internal/ui"
	"github.com/fraol163/viren/internal/util"
)

const (
	GitHubAPIURL	= "https://api.github.com/repos/fraol163/viren/releases/latest"
	GitHubRepo	= "fraol163/viren"
	UpdateTimeout	= 300 * time.Second
	ChunkSize	= 1024 * 1024 // 1MB chunks for resume capability
)

type ReleaseInfo struct {
	TagName	string		`json:"tag_name"`
	Name	string		`json:"name"`
	Body	string		`json:"body"`
	PublishedAt	string		`json:"published_at"`
	Assets	[]Asset		`json:"assets"`
}

type Asset struct {
	Name		string	`json:"name"`
	DownloadURL	string	`json:"browser_download_url"`
	Size		int64	`json:"size"`
}

type UpdateManager struct {
	currentVersion	string
	terminal	*ui.Terminal
}

func NewManager(currentVersion string, terminal *ui.Terminal) *UpdateManager {
	return &UpdateManager{
		currentVersion: currentVersion,
		terminal:	terminal,
	}
}

func (u *UpdateManager) CheckForUpdates() (bool, *ReleaseInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", GitHubAPIURL, nil)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "viren-updater")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(u.currentVersion, "v")

	hasUpdate := latestVersion != currentVersion
	return hasUpdate, &release, nil
}

func (u *UpdateManager) DownloadUpdate(release *ReleaseInfo, cancel *context.CancelFunc) error {
	assetName := u.getAssetName()
	var downloadURL string
	var fileSize int64

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.DownloadURL
			fileSize = asset.Size
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no suitable binary found for %s/%s (looking for %s)", runtime.GOOS, runtime.GOARCH, assetName)
	}

	u.terminal.PrintInfo(fmt.Sprintf("downloading viren %s (%s)...", strings.TrimPrefix(release.TagName, "v"), formatSize(fileSize)))

	tempDir, err := util.GetTempDir()
	if err != nil {
		return fmt.Errorf("failed to get temp directory: %w", err)
	}

	partialFile := filepath.Join(tempDir, "viren-update-"+release.TagName+".partial")

	return u.downloadWithResume(downloadURL, partialFile, fileSize, cancel)
}

func (u *UpdateManager) downloadWithResume(url, partialFile string, fileSize int64, cancel *context.CancelFunc) error {
	var startOffset int64 = 0

	if info, err := os.Stat(partialFile); err == nil {
		startOffset = info.Size()
		if startOffset >= fileSize {
			u.terminal.PrintInfo("download already complete, verifying...")
			return nil
		}
		u.terminal.PrintInfo(fmt.Sprintf("resuming download from %s...", formatSize(startOffset)))
	}

	file, err := os.OpenFile(partialFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ctx, cancelCtx := context.WithTimeout(context.Background(), UpdateTimeout)
	if cancel != nil {
		defer cancelCtx()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if startOffset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startOffset))
	}
	req.Header.Set("User-Agent", "viren-updater")

	client := &http.Client{Timeout: UpdateTimeout}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("download timeout - will resume on next run")
		}
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	totalDownloaded := startOffset
	lastProgress := time.Now()
	progressWidth := 50

	fmt.Println()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("download interrupted - will resume on next run")
		default:
		}

		buffer := make([]byte, ChunkSize)
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := file.Write(buffer[:n])
			if writeErr != nil {
				return fmt.Errorf("failed to write: %w", writeErr)
			}
			totalDownloaded += int64(n)

			if time.Since(lastProgress) > 500*time.Millisecond {
				progress := float64(totalDownloaded) / float64(fileSize) * 100
				filled := int(progress / 100 * float64(progressWidth))
				bar := strings.Repeat("█", filled) + strings.Repeat("░", progressWidth-filled)
				fmt.Printf("\r\033[K[%s] %.1f%% %s/%s", bar, progress, formatSize(totalDownloaded), formatSize(fileSize))
				lastProgress = time.Now()
			}
		}

		if err == io.EOF {
			fmt.Println()
			break
		}
		if err != nil {
			return fmt.Errorf("download error: %w", err)
		}
	}

	fmt.Println()
	u.terminal.PrintInfo("download complete, installing...")
	return nil
}

func (u *UpdateManager) InstallUpdate(release *ReleaseInfo) error {
	tempDir, err := util.GetTempDir()
	if err != nil {
		return fmt.Errorf("failed to get temp directory: %w", err)
	}

	partialFile := filepath.Join(tempDir, "viren-update-"+release.TagName+".partial")
	newBinary := filepath.Join(tempDir, "viren-new-"+release.TagName)

	if _, err := os.Stat(partialFile); os.IsNotExist(err) {
		return fmt.Errorf("downloaded file not found")
	}

	if err := os.Rename(partialFile, newBinary); err != nil {
		return fmt.Errorf("failed to prepare binary: %w", err)
	}

	if err := os.Chmod(newBinary, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	currentExec, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	backupFile := currentExec + ".backup"
	if err := copyFile(currentExec, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	if err := os.Rename(newBinary, currentExec); err != nil {
		u.terminal.PrintError("update failed, restoring backup...")
		if restoreErr := os.Rename(backupFile, currentExec); restoreErr != nil {
			return fmt.Errorf("failed to restore backup: %v (original error: %v)", restoreErr, err)
		}
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	os.Remove(backupFile)
	os.Remove(filepath.Join(tempDir, "viren-update-"+release.TagName))

	u.terminal.PrintInfo("update successful!")
	time.Sleep(1 * time.Second)

	u.displayWhatsNew(release)

	time.Sleep(3 * time.Second)
	u.terminal.PrintInfo("restarting...")
	time.Sleep(2 * time.Second)

	cmd := exec.Command(currentExec, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to restart: %w", err)
	}

	os.Exit(0)
	return nil
}

func (u *UpdateManager) displayWhatsNew(release *ReleaseInfo) {
	theme := u.terminal.GetTheme()
	
	fmt.Println()
	fmt.Printf("%s WHAT'S NEW \033[0m\n", theme.SuccessBox)
	fmt.Printf("%s\n", strings.Repeat("═", 70))
	
	if release.Body != "" {
		lines := strings.Split(release.Body, "\n")
		inFeatures := false
		
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			
			if strings.Contains(strings.ToLower(trimmed), "feature") || 
			   strings.Contains(strings.ToLower(trimmed), "change") ||
			   strings.Contains(strings.ToLower(trimmed), "improvement") ||
			   strings.HasPrefix(trimmed, "-") ||
			   strings.HasPrefix(trimmed, "•") ||
			   strings.HasPrefix(trimmed, "*") ||
			   strings.HasPrefix(trimmed, "✅") {
				inFeatures = true
			}
			
			if inFeatures && trimmed != "" {
				if strings.HasPrefix(trimmed, "-") || 
				   strings.HasPrefix(trimmed, "•") || 
				   strings.HasPrefix(trimmed, "*") ||
				   strings.HasPrefix(trimmed, "✅") {
					fmt.Printf("  \033[92m▸\033[0m %s\n", strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(trimmed, "-"), "•"), "*"))
				} else if !strings.HasPrefix(trimmed, "#") && len(trimmed) > 3 {
					fmt.Printf("  \033[96m%s\033[0m\n", trimmed)
				}
			}
		}
	}
	
	fmt.Printf("%s\n", strings.Repeat("═", 70))
	fmt.Println()
}

func (u *UpdateManager) getAssetName() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goos == "windows" {
		return fmt.Sprintf("viren_%s_%s.exe", goos, goarch)
	}
	return fmt.Sprintf("viren_%s_%s", goos, goarch)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (u *UpdateManager) CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (u *UpdateManager) Cleanup() {
	tempDir, err := util.GetTempDir()
	if err != nil {
		return
	}

	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasPrefix(info.Name(), "viren-update-") || strings.HasPrefix(info.Name(), "viren-new-") {
			if info.ModTime().Add(24 * time.Hour).Before(time.Now()) {
				os.Remove(path)
			}
		}
		return nil
	})
}
