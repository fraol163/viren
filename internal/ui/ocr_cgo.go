//go:build cgo

package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

func (t *Terminal) extractTextFromImage(filePath string) (string, error) {

	if _, err := exec.LookPath("tesseract"); err != nil {
		return "", fmt.Errorf("tesseract OCR is not installed. Please install it to enable image-to-text extraction")
	}

	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetLanguage("eng")
	if err != nil {

		client = gosseract.NewClient()
		defer client.Close()
	}

	err = client.SetImage(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to set image source: %w", err)
	}

	client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.,!?@#$%^&*()_+-={}[]|\\:;\"'<>/~` ")

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("OCR extraction failed: %w", err)
	}

	cleanedText := strings.TrimSpace(text)

	lines := strings.Split(cleanedText, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n"), nil
}
