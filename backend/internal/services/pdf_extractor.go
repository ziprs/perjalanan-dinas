package services

import (
	"fmt"
	"os/exec"
	"strings"
)

// PDFExtractor handles PDF text extraction
type PDFExtractor struct{}

// NewPDFExtractor creates a new PDF extractor
func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

// ExtractText extracts text from PDF file using pdftotext
func (pe *PDFExtractor) ExtractText(pdfPath string) (string, error) {
	// Try pdftotext first (most accurate)
	if text, err := pe.extractWithPdfToText(pdfPath); err == nil && text != "" {
		return text, nil
	}

	// Fallback to strings command
	return pe.extractWithStrings(pdfPath)
}

// extractWithPdfToText uses pdftotext utility
func (pe *PDFExtractor) extractWithPdfToText(pdfPath string) (string, error) {
	cmd := exec.Command("pdftotext", "-layout", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// extractWithStrings uses strings command as fallback
func (pe *PDFExtractor) extractWithStrings(pdfPath string) (string, error) {
	cmd := exec.Command("strings", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// ExtractTextClean extracts and cleans text from PDF
func (pe *PDFExtractor) ExtractTextClean(pdfPath string) (string, error) {
	text, err := pe.ExtractText(pdfPath)
	if err != nil {
		return "", err
	}

	// Clean up text
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Remove excessive blank lines
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	previousBlank := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !previousBlank {
				cleanedLines = append(cleanedLines, "")
				previousBlank = true
			}
		} else {
			cleanedLines = append(cleanedLines, trimmed)
			previousBlank = false
		}
	}

	return strings.Join(cleanedLines, "\n"), nil
}

// ValidatePDF checks if file is a valid PDF
func (pe *PDFExtractor) ValidatePDF(pdfPath string) error {
	cmd := exec.Command("file", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check file type: %w", err)
	}

	outputStr := strings.ToLower(string(output))
	if !strings.Contains(outputStr, "pdf") {
		return fmt.Errorf("file is not a PDF document")
	}

	return nil
}
