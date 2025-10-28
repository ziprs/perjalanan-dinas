package utils

import (
	"fmt"
	"regexp"
)

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

// ExtractPositionCodeFromRequestNumber extracts position code from request number
// Example: "064/0325/DIB/DPEB/NOTA" -> "DPEB"
func ExtractPositionCodeFromRequestNumber(requestNumber string) string {
	re := regexp.MustCompile(`/DIB/([A-Z]+)/NOTA`)
	matches := re.FindStringSubmatch(requestNumber)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
