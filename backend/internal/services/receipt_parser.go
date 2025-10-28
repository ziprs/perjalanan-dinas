package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ReceiptData represents parsed data from receipt
type ReceiptData struct {
	ReceiptNumber   string    `json:"receipt_number"`
	ReceiptDate     time.Time `json:"receipt_date"`
	Vendor          string    `json:"vendor"`
	PassengerName   string    `json:"passenger_name"`
	Type            string    `json:"type"` // flight, hotel, train
	Description     string    `json:"description"`
	Amount          int       `json:"amount"`
	RouteOrLocation string    `json:"route_or_location"`
	RawText         string    `json:"raw_text,omitempty"`
}

// ReceiptParser handles receipt parsing
type ReceiptParser struct{}

// NewReceiptParser creates a new receipt parser
func NewReceiptParser() *ReceiptParser {
	return &ReceiptParser{}
}

// ParseReceiptText parses text content from receipt PDF
func (rp *ReceiptParser) ParseReceiptText(text string) (*ReceiptData, error) {
	data := &ReceiptData{
		RawText: text,
	}

	// Detect vendor (Traveloka, tiket.com, etc)
	data.Vendor = rp.detectVendor(text)

	// Parse receipt number
	data.ReceiptNumber = rp.extractReceiptNumber(text)

	// Parse date
	if date, err := rp.extractDate(text); err == nil {
		data.ReceiptDate = date
	}

	// Parse passenger name
	data.PassengerName = rp.extractPassengerName(text)

	// Detect type (flight, hotel, train)
	data.Type = rp.detectType(text)

	// Parse amount
	data.Amount = rp.extractAmount(text)

	// Parse route or location
	data.RouteOrLocation = rp.extractRouteOrLocation(text, data.Type)

	// Generate description
	data.Description = rp.generateDescription(data)

	return data, nil
}

// detectVendor detects the travel vendor from text
func (rp *ReceiptParser) detectVendor(text string) string {
	textLower := strings.ToLower(text)

	if strings.Contains(textLower, "traveloka") || strings.Contains(textLower, "trinusa travelindo") {
		return "Traveloka"
	}
	if strings.Contains(textLower, "tiket.com") || strings.Contains(textLower, "tiket com") {
		return "Tiket.com"
	}
	if strings.Contains(textLower, "pegipegi") {
		return "PegiPegi"
	}
	// Detect travel agencies / biro travel
	if strings.Contains(textLower, "sonyloka") || strings.Contains(textLower, "tiketqta") {
		return "Sonyloka Travel"
	}
	if strings.Contains(textLower, "biro") && strings.Contains(textLower, "travel") {
		return "Biro Travel"
	}
	if strings.Contains(textLower, "travel agent") || strings.Contains(textLower, "tour") {
		return "Travel Agent"
	}
	if strings.Contains(textLower, "receipt") || strings.Contains(textLower, "bukti pembelian") {
		return "Travel Provider"
	}

	return "Unknown"
}

// extractReceiptNumber extracts receipt/booking number
func (rp *ReceiptParser) extractReceiptNumber(text string) string {
	// Pattern: #1844870088383028819 or similar
	patterns := []string{
		`#(\d{13,25})`,
		`Nomor\s*:\s*#(\d{13,25})`,
		`Number\s*:\s*#(\d{13,25})`,
		`Receipt.*#(\d{13,25})`,
		// Travel agency patterns
		`KODE RESERVASI VOUCHER HOTEL:\s*(\d+)`,
		`Kode Booking:\s*(\w+)`,
		`Booking Code:\s*(\w+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			// Add # prefix only for numeric-only codes without it
			if regexp.MustCompile(`^\d+$`).MatchString(matches[1]) && !strings.HasPrefix(matches[1], "#") {
				return "#" + matches[1]
			}
			return matches[1]
		}
	}

	return ""
}

// extractDate extracts receipt date
func (rp *ReceiptParser) extractDate(text string) (time.Time, error) {
	// Pattern: 02 Okt 2025, 23:41 or 01 Oct 2025, 09:19 or 17/10/2021
	patterns := []string{
		`(\d{2})\s+(Okt?|Oct)\s+(\d{4})`,
		`(\d{2})\s+(Sep?t?|Nov|Des?c?|Jan|Feb|Mar|Apr|Mei?|May|Jun|Jul|Agu?s?|Aug)\s+(\d{4})`,
		// Travel agency DD/MM/YYYY format
		`(\d{2})/(\d{2})/(\d{4})`,
	}

	monthMap := map[string]string{
		"jan": "Jan", "feb": "Feb", "mar": "Mar", "apr": "Apr",
		"mei": "May", "may": "May", "jun": "Jun", "jul": "Jul",
		"agu": "Aug", "aug": "Aug", "sep": "Sep", "sept": "Sep",
		"okt": "Oct", "oct": "Oct", "nov": "Nov", "des": "Dec", "dec": "Dec",
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 3 {
			// Check if it's DD/MM/YYYY format
			if strings.Contains(pattern, "/") {
				day := matches[1]
				month := matches[2]
				year := matches[3]
				dateStr := fmt.Sprintf("%s/%s/%s", day, month, year)
				if t, err := time.Parse("02/01/2006", dateStr); err == nil {
					return t, nil
				}
			} else {
				day := matches[1]
				month := strings.ToLower(matches[2])
				year := matches[3]

				// Convert Indonesian month to English
				if englishMonth, ok := monthMap[month[:3]]; ok {
					month = englishMonth
				}

				dateStr := fmt.Sprintf("%s %s %s", day, month, year)
				formats := []string{
					"02 Jan 2006",
					"02 January 2006",
				}

				for _, format := range formats {
					if t, err := time.Parse(format, dateStr); err == nil {
						return t, nil
					}
				}
			}
		}
	}

	return time.Now(), fmt.Errorf("date not found")
}

// extractPassengerName extracts passenger/guest name
func (rp *ReceiptParser) extractPassengerName(text string) string {
	// Pattern for flight: "Tn. Erfan Basrianto (DEWASA)"
	// Pattern for hotel: "Erfan Basrianto"

	patterns := []string{
		`(?:Tn\.|Mr\.|Ny\.|Ms\.)\s+([A-Za-z\s]+)\s+\((?:DEWASA|ADULT|Dewasa|Adult)\)`,
		`(?:Nama|Name)\s*:\s*([A-Za-z\s]+)`,
		`TAMU[\s\n]+([A-Za-z\s]+)`,
		`PASSENGER DETAILS[\s\n]+(?:Mr\.|Ms\.|Tn\.|Ny\.)\s+([A-Za-z\s]+)`,
		// Travel agency patterns
		`Data Tamu Pemesan[\s\S]*?\d+\.\s*(?:Mr\.|Ms\.|Tn\.|Ny\.)\s*([A-Za-z\s]+)`,
		`Nama Pemesan[:\s]+([A-Za-z\s]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			name := strings.TrimSpace(matches[1])
			// Clean up name (remove extra spaces)
			name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
			return name
		}
	}

	return ""
}

// detectType detects the type of booking
func (rp *ReceiptParser) detectType(text string) string {
	textLower := strings.ToLower(text)

	// Check for flight indicators
	if strings.Contains(textLower, "tiket pesawat") ||
	   strings.Contains(textLower, "flight ticket") ||
	   strings.Contains(textLower, "batik air") ||
	   strings.Contains(textLower, "garuda") ||
	   strings.Contains(textLower, "lion air") ||
	   regexp.MustCompile(`[A-Z]{3}\s*-\s*[A-Z]{3}`).MatchString(text) {
		return "flight"
	}

	// Check for hotel indicators
	if strings.Contains(textLower, "hotel") ||
	   strings.Contains(textLower, "akomodasi") ||
	   strings.Contains(textLower, "accommodation") ||
	   strings.Contains(textLower, "check-in") ||
	   strings.Contains(textLower, "kamar") {
		return "hotel"
	}

	// Check for train indicators
	if strings.Contains(textLower, "kereta") ||
	   strings.Contains(textLower, "train") ||
	   strings.Contains(textLower, "kai") {
		return "train"
	}

	return "other"
}

// extractAmount extracts total payment amount
func (rp *ReceiptParser) extractAmount(text string) int {
	// Look for "JUMLAH PEMBAYARAN" or "PAYMENT AMOUNT" or "TOTAL"
	patterns := []string{
		`JUMLAH PEMBAYARAN\s+(?:Rp\s*)?([\d\.]+)`,
		`PAYMENT AMOUNT\s+(?:Rp\s*)?([\d\.]+)`,
		`TOTAL\s+(?:Rp\s*)?([\d\.]+)`,
		`Total Rp\s+([\d\.]+)`,
		// Travel agency patterns with "HARGA" or "Rp." (with period)
		`HARGA\s+Rp\.\s*([\d\.]+)`,
		// Fallback: any "Rp" or "Rp." followed by numbers with dots
		`(?:Rp\.?|IDR)\s*([\d\.]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(text); len(matches) > 1 {
			// Remove dots (thousands separator in Indonesian format)
			amountStr := strings.ReplaceAll(matches[1], ".", "")
			// Remove commas if any
			amountStr = strings.ReplaceAll(amountStr, ",", "")
			// Try to convert to int
			if amount, err := strconv.Atoi(amountStr); err == nil && amount > 0 {
				return amount
			}
		}
	}

	return 0
}

// extractRouteOrLocation extracts flight route or hotel location
func (rp *ReceiptParser) extractRouteOrLocation(text string, receiptType string) string {
	if receiptType == "flight" {
		// Pattern: HLP - SUB or SUB - HLP
		re := regexp.MustCompile(`([A-Z]{3})\s*-\s*([A-Z]{3})`)
		if matches := re.FindStringSubmatch(text); len(matches) > 2 {
			return fmt.Sprintf("%s - %s", matches[1], matches[2])
		}

		// Pattern: Batik Air (Dewasa) HLP - SUB
		re = regexp.MustCompile(`(?:Batik Air|Garuda|Lion Air).*?([A-Z]{3})\s*-\s*([A-Z]{3})`)
		if matches := re.FindStringSubmatch(text); len(matches) > 2 {
			return fmt.Sprintf("%s - %s", matches[1], matches[2])
		}
	} else if receiptType == "hotel" {
		// Extract hotel name - stop at newline or certain keywords
		patterns := []string{
			`DETAIL HOTEL[\s\n]+([A-Za-z\s]+?)(?:\n|Alamat|Check)`,
			`Akomodasi\s+([A-Za-z\s]+?)\s*,`,
			// Travel agency pattern: "Detail Booking HOTEL NAME - (date)"
			`Detail Booking\s+([A-Za-z\s]+?)\s+(?:Hotel|by|-)`,
			`([A-Za-z\s]+Hotel)`,
		}

		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			if matches := re.FindStringSubmatch(text); len(matches) > 1 {
				hotelName := strings.TrimSpace(matches[1])
				// Additional cleanup - remove extra whitespace
				hotelName = regexp.MustCompile(`\s+`).ReplaceAllString(hotelName, " ")
				return hotelName
			}
		}
	}

	return ""
}

// generateDescription generates a human-readable description
func (rp *ReceiptParser) generateDescription(data *ReceiptData) string {
	switch data.Type {
	case "flight":
		if data.RouteOrLocation != "" {
			return fmt.Sprintf("Tiket Pesawat %s", data.RouteOrLocation)
		}
		return "Tiket Pesawat"
	case "hotel":
		if data.RouteOrLocation != "" {
			return fmt.Sprintf("Hotel %s", data.RouteOrLocation)
		}
		return "Akomodasi Hotel"
	case "train":
		if data.RouteOrLocation != "" {
			return fmt.Sprintf("Tiket Kereta %s", data.RouteOrLocation)
		}
		return "Tiket Kereta"
	default:
		return "Akomodasi Perjalanan Dinas"
	}
}

// ToJSON converts parsed data to JSON string
func (rp *ReceiptParser) ToJSON(data *ReceiptData) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// FromJSON parses JSON string to ReceiptData
func (rp *ReceiptParser) FromJSON(jsonStr string) (*ReceiptData, error) {
	var data ReceiptData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
