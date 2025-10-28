package services

import (
	"testing"
)

// Test sample receipt text from Traveloka
const sampleTravelokaFlight = `
BUKTI PEMBELIAN (RECEIPT)
Nomor : #1844870088383028819
Tanggal : 02 Okt 2025, 23:41 (Kamis)

DATA PEMESAN
Nama : Erfan Basrianto
Email : ebastian9987@gmail.com

DATA PENUMPANG
Tn. Erfan Basrianto (DEWASA) | Tn. Erfan Basrianto (Dewasa)

DETAIL PEMBELIAN
No. Jenis Barang Deskripsi Jml. Harga satuan Rp Total Rp
1 Tiket Pesawat Batik Air (Dewasa) HLP - SUB | 3 Okt 2025 1 1.540.600 1.540.600
2 Asuransi Perjalanan Asuransi Perjalanan - Plan Domestik 1 37.000 37.000

TOTAL 1.577.600
BIAYA LAYANAN * 22.000
JUMLAH PEMBAYARAN 1.599.600
`

const sampleTravelokaHotel = `
BUKTI PEMBELIAN (RECEIPT)
Nomor : #1844827144826419713
Tanggal : 01 Okt 2025, 09:19 (Rabu)

DATA PEMESAN
Nama : Erfan Basrianto

TAMU
Erfan Basrianto

DETAIL HOTEL
Ashley Tanah Abang
Alamat: Jalan K.H Wahid Hasyim 220
Check-in: 01-10-2025
Durasi: 2 malam

DETAIL PEMBELIAN
1 Akomodasi Ashley Tanah Abang , Kamar Double Superior - 2 tamu 1 1.673.558 1.673.558
2 Pajak dan Biaya Lainnya Pajak dan Biaya Lainnya 1 225.930 225.930

TOTAL 1.899.488
JUMLAH PEMBAYARAN 1.899.488
`

func TestParseFlightReceipt(t *testing.T) {
	parser := NewReceiptParser()
	data, err := parser.ParseReceiptText(sampleTravelokaFlight)

	if err != nil {
		t.Fatalf("Failed to parse receipt: %v", err)
	}

	// Test vendor detection
	if data.Vendor != "Traveloka" {
		t.Errorf("Expected vendor 'Traveloka', got '%s'", data.Vendor)
	}

	// Test receipt number
	if data.ReceiptNumber != "#1844870088383028819" {
		t.Errorf("Expected receipt number '#1844870088383028819', got '%s'", data.ReceiptNumber)
	}

	// Test passenger name
	if data.PassengerName != "Erfan Basrianto" {
		t.Errorf("Expected passenger 'Erfan Basrianto', got '%s'", data.PassengerName)
	}

	// Test type
	if data.Type != "flight" {
		t.Errorf("Expected type 'flight', got '%s'", data.Type)
	}

	// Test amount
	if data.Amount != 1599600 {
		t.Errorf("Expected amount 1599600, got %d", data.Amount)
	}

	// Test route
	if data.RouteOrLocation != "HLP - SUB" {
		t.Errorf("Expected route 'HLP - SUB', got '%s'", data.RouteOrLocation)
	}

	t.Logf("✓ Flight receipt parsed successfully:")
	t.Logf("  Vendor: %s", data.Vendor)
	t.Logf("  Receipt: %s", data.ReceiptNumber)
	t.Logf("  Passenger: %s", data.PassengerName)
	t.Logf("  Type: %s", data.Type)
	t.Logf("  Amount: Rp %d", data.Amount)
	t.Logf("  Route: %s", data.RouteOrLocation)
	t.Logf("  Description: %s", data.Description)
}

func TestParseHotelReceipt(t *testing.T) {
	parser := NewReceiptParser()
	data, err := parser.ParseReceiptText(sampleTravelokaHotel)

	if err != nil {
		t.Fatalf("Failed to parse receipt: %v", err)
	}

	// Test vendor detection
	if data.Vendor != "Traveloka" {
		t.Errorf("Expected vendor 'Traveloka', got '%s'", data.Vendor)
	}

	// Test type
	if data.Type != "hotel" {
		t.Errorf("Expected type 'hotel', got '%s'", data.Type)
	}

	// Test amount
	if data.Amount != 1899488 {
		t.Errorf("Expected amount 1899488, got %d", data.Amount)
	}

	// Test location (hotel name)
	if data.RouteOrLocation != "Ashley Tanah Abang" {
		t.Errorf("Expected location 'Ashley Tanah Abang', got '%s'", data.RouteOrLocation)
	}

	t.Logf("✓ Hotel receipt parsed successfully:")
	t.Logf("  Vendor: %s", data.Vendor)
	t.Logf("  Receipt: %s", data.ReceiptNumber)
	t.Logf("  Guest: %s", data.PassengerName)
	t.Logf("  Type: %s", data.Type)
	t.Logf("  Amount: Rp %d", data.Amount)
	t.Logf("  Hotel: %s", data.RouteOrLocation)
	t.Logf("  Description: %s", data.Description)
}

func TestJSONConversion(t *testing.T) {
	parser := NewReceiptParser()
	data, _ := parser.ParseReceiptText(sampleTravelokaFlight)

	// Convert to JSON
	jsonStr, err := parser.ToJSON(data)
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	// Parse back from JSON
	parsedData, err := parser.FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse from JSON: %v", err)
	}

	// Verify data integrity
	if parsedData.ReceiptNumber != data.ReceiptNumber {
		t.Errorf("JSON conversion failed: receipt number mismatch")
	}

	t.Logf("✓ JSON conversion successful")
}
