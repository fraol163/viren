//go:build !cgo

package ui

import "fmt"

func (t *Terminal) extractTextFromImage(filePath string) (string, error) {
	return "", fmt.Errorf("OCR (Tesseract) is not available on this platform. Image-to-text extraction is disabled")
}
