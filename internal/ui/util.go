package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fraol163/viren/pkg/types"
)

// GetAvailableEditors returns a list of installed editors based on the operating system
func GetAvailableEditors(cfg *types.Config) []string {
	var potentialEditors []string

	if cfg.PreferredEditor != "" {
		potentialEditors = append(potentialEditors, cfg.PreferredEditor)
	}

	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		potentialEditors = append(potentialEditors, envEditor)
	}

	if runtime.GOOS == "windows" {
		potentialEditors = append(potentialEditors, "notepad", "code", "notepad++")
	} else if runtime.GOOS == "linux" {
		potentialEditors = append(potentialEditors, "gedit", "mousepad", "kate", "leafpad", "xed", "code")
	} else if runtime.GOOS == "darwin" {
		potentialEditors = append(potentialEditors, "TextEdit", "code")
	}

	potentialEditors = append(potentialEditors, "vim", "vi", "nano", "nvim", "emacs")

	seen := make(map[string]bool)
	var available []string
	for _, editor := range potentialEditors {
		if editor == "" || seen[editor] {
			continue
		}
		if _, err := exec.LookPath(editor); err == nil {
			available = append(available, editor)
			seen[editor] = true
		}
	}

	return available
}

// RunSpecificEditor runs a chosen editor on a file
func RunSpecificEditor(editor string, filePath string) error {
	var cmd *exec.Cmd

	if editor == "TextEdit" && runtime.GOOS == "darwin" {
		cmd = exec.Command("open", "-e", filePath)
	} else {
		cmd = exec.Command(editor, filePath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RunEditorWithFallback tries to run the user's preferred editor, then falls back to common editors.
func RunEditorWithFallback(cfg *types.Config, filePath string) error {
	editors := GetAvailableEditors(cfg)
	if len(editors) == 0 {
		return fmt.Errorf("no working editor found")
	}

	for _, editor := range editors {
		if err := RunSpecificEditor(editor, filePath); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to launch any editor")
}
