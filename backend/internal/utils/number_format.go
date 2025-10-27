package utils

import "fmt"

// FormatSequence formats integer to padded string
func FormatSequence(num int) string {
	return fmt.Sprintf("%04d", num)
}

// GenerateRequestNumber formats request number: 064/{seq}/DIB/{code}/NOTA
func GenerateRequestNumber(seq int, positionCode string) string {
	return fmt.Sprintf("064/%s/DIB/%s/NOTA", FormatSequence(seq), positionCode)
}

// GenerateReportNumber formats report number: 064/{seq}/DIB/{code}/NOTA
func GenerateReportNumber(seq int, positionCode string) string {
	return fmt.Sprintf("064/%s/DIB/%s/NOTA", FormatSequence(seq), positionCode)
}
