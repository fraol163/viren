package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fraol163/viren/pkg/types"
)

func GetTempDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	tempDir := filepath.Join(homeDir, ".viren", "tmp")

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return tempDir, nil
}

func IsShallowLoadDir(cfg *types.Config, dirPath string) bool {

	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return false
	}
	absPath = filepath.Clean(absPath)

	for _, shallowDir := range cfg.ShallowLoadDirs {
		if shallowDir == "" {
			continue
		}

		if strings.HasPrefix(shallowDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				continue
			}
			shallowDir = filepath.Join(homeDir, shallowDir[1:])
		}

		absShallowDir, err := filepath.Abs(shallowDir)
		if err != nil {
			continue
		}
		absShallowDir = filepath.Clean(absShallowDir)

		if absPath == absShallowDir {
			return true
		}
	}

	return false
}
