package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents admin user
type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Position represents jabatan with allowance rates
type Position struct {
	ID                       uint           `gorm:"primarykey" json:"id"`
	Title                    string         `gorm:"unique;not null" json:"title"`                       // Nama Jabatan
	Code                     string         `gorm:"not null" json:"code"`                               // Kode (DPEB, DPDB, DDBE)
	Level                    string         `gorm:"not null" json:"level"`                              // Jr. Officer, Officer, Senior Officer, AVP
	AllowanceInProvince      int            `gorm:"not null" json:"allowance_in_province"`              // Tarif dalam provinsi
	AllowanceOutsideProvince int            `gorm:"not null" json:"allowance_outside_province"`         // Tarif luar provinsi
	AllowanceAbroad          int            `gorm:"not null" json:"allowance_abroad"`                   // Tarif luar negeri
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// Employee represents karyawan
type Employee struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	NIP        string         `gorm:"column:nip;unique;not null" json:"nip"`
	Name       string         `gorm:"not null" json:"name"`
	PositionID uint           `gorm:"not null" json:"position_id"` // Foreign key ke positions
	Position   Position       `gorm:"foreignKey:PositionID" json:"position"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TravelRequest represents nota permintaan surat tugas perjalanan dinas
type TravelRequest struct {
	ID                     uint                    `gorm:"primarykey" json:"id"`
	Purpose                string                  `gorm:"type:text;not null" json:"purpose"`                     // Maksud perjalanan dinas
	DeparturePlace         string                  `gorm:"not null;default:'Surabaya'" json:"departure_place"`    // Tempat berangkat asal
	Destination            string                  `gorm:"not null" json:"destination"`                           // Tempat tujuan (nama kota)
	DestinationType        string                  `gorm:"not null" json:"destination_type"`                      // "in_province", "outside_province", "abroad"
	DepartureDate          time.Time               `gorm:"not null" json:"departure_date"`
	ReturnDate             time.Time               `gorm:"not null" json:"return_date"`
	DurationDays           int                     `gorm:"not null" json:"duration_days"`                         // Lama perjalanan dinas (auto calculated)
	Transportation         string                  `gorm:"not null" json:"transportation"`                        // angkutan umum, pesawat, kereta api
	TotalAllowance         int                     `gorm:"not null;default:0" json:"total_allowance"`             // Total iuran (auto calculated)
	RequestNumber          string                  `gorm:"unique;not null" json:"request_number"`                 // 064/{seq}/DIB/{code}/NOTA
	ReportNumber           string                  `gorm:"unique" json:"report_number"`                           // 064/ /DIB/{code}/NOTA
	Status                 string                  `gorm:"default:'pending'" json:"status"`                       // pending, approved, completed
	TravelRequestEmployees []TravelRequestEmployee `gorm:"foreignKey:TravelRequestID" json:"employees"`
	TravelReport           *TravelReport           `gorm:"foreignKey:TravelRequestID" json:"travel_report,omitempty"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
	DeletedAt              gorm.DeletedAt          `gorm:"index" json:"-"`
}

// TravelRequestEmployee represents the many-to-many relationship between TravelRequest and Employee
type TravelRequestEmployee struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	TravelRequestID  uint           `gorm:"not null" json:"travel_request_id"`
	EmployeeID       uint           `gorm:"not null" json:"employee_id"`
	Employee         Employee       `gorm:"foreignKey:EmployeeID" json:"employee"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TravelReport represents berita acara perjalanan dinas
type TravelReport struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	TravelRequestID   uint           `gorm:"unique;not null" json:"travel_request_id"`
	TravelRequest     TravelRequest  `gorm:"foreignKey:TravelRequestID;constraint:OnDelete:CASCADE" json:"travel_request"`
	ReportNumber      string         `gorm:"unique;not null" json:"report_number"` // 064/ /DIB/{code}/NOTA
	RepresentativeName string        `gorm:"not null" json:"representative_name"` // Nama perwakilan yang tanda tangan
	RepresentativePosition string    `gorm:"not null" json:"representative_position"` // Jabatan perwakilan
	VisitProofs       []VisitProof   `gorm:"foreignKey:TravelReportID" json:"visit_proofs"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// VisitProof represents bukti kunjungan (tabel pada halaman 2)
type VisitProof struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	TravelReportID  uint           `gorm:"not null" json:"travel_report_id"`
	Date            time.Time      `gorm:"not null" json:"date"`
	DepartFrom      string         `gorm:"not null" json:"depart_from"`      // Berangkat dari
	StayOrStopAt    string         `json:"stay_or_stop_at"`                  // Bermalam atau singgah di
	ArriveAt        string         `gorm:"not null" json:"arrive_at"`        // Datang di
	SignatureProof  string         `json:"signature_proof"`                  // Path/info tanda tangan & stempel
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// NumberingConfig untuk manage penomoran
type NumberingConfig struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	LastRequestSequence   int            `gorm:"not null;default:0" json:"last_request_sequence"`
	LastReportSequence    int            `gorm:"not null;default:0" json:"last_report_sequence"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

// RepresentativeConfig untuk manage perwakilan penandatangan
type RepresentativeConfig struct {
	ID       uint           `gorm:"primarykey" json:"id"`
	Name     string         `gorm:"not null" json:"name"`
	Position string         `gorm:"not null" json:"position"`
	IsActive bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AtCostClaim represents claim untuk biaya transportasi dan akomodasi
type AtCostClaim struct {
	ID                     uint                    `gorm:"primarykey" json:"id"`
	TravelRequestID        uint                    `gorm:"not null" json:"travel_request_id"`
	TravelRequest          TravelRequest           `gorm:"foreignKey:TravelRequestID" json:"travel_request"`
	ClaimNumber            string                  `gorm:"unique;not null" json:"claim_number"`                 // 064/{seq}/DIB/{code}/NOTA (same sequence)
	RepresentativeName     string                  `gorm:"not null" json:"representative_name"`                 // VP name
	RepresentativePosition string                  `gorm:"not null" json:"representative_position"`             // VP position
	Status                 string                  `gorm:"default:'pending'" json:"status"`                     // pending, approved, rejected
	TotalAmount            int                     `gorm:"not null;default:0" json:"total_amount"`              // Total semua klaim
	ClaimItems             []AtCostClaimItem       `gorm:"foreignKey:AtCostClaimID" json:"claim_items"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
	DeletedAt              gorm.DeletedAt          `gorm:"index" json:"-"`
}

// AtCostClaimItem represents detail klaim per karyawan
type AtCostClaimItem struct {
	ID              uint                      `gorm:"primarykey" json:"id"`
	AtCostClaimID   uint                      `gorm:"not null" json:"at_cost_claim_id"`
	EmployeeID      uint                      `gorm:"not null" json:"employee_id"`
	Employee        Employee                  `gorm:"foreignKey:EmployeeID" json:"employee"`
	TransportCost   int                       `gorm:"not null;default:0" json:"transport_cost"`    // Biaya transportasi
	AccommodationCost int                     `gorm:"not null;default:0" json:"accommodation_cost"` // Biaya penginapan
	TotalCost       int                       `gorm:"not null;default:0" json:"total_cost"`        // Total per karyawan
	Receipts        []AtCostReceipt           `gorm:"foreignKey:ClaimItemID" json:"receipts"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	DeletedAt       gorm.DeletedAt            `gorm:"index" json:"-"`
}

// AtCostReceipt represents bukti pembayaran (invoice/receipt)
type AtCostReceipt struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	ClaimItemID     uint           `gorm:"not null" json:"claim_item_id"`
	ReceiptNumber   string         `json:"receipt_number"`                        // Nomor receipt dari Traveloka/tiket.com
	ReceiptDate     time.Time      `json:"receipt_date"`                          // Tanggal receipt
	Vendor          string         `json:"vendor"`                                // Traveloka, tiket.com, etc
	Type            string         `gorm:"not null" json:"type"`                  // flight, hotel, train
	Description     string         `gorm:"type:text" json:"description"`          // Deskripsi item
	Amount          int            `gorm:"not null" json:"amount"`                // Nominal
	FilePath        string         `gorm:"not null" json:"file_path"`             // Path ke file PDF
	FileName        string         `gorm:"not null" json:"file_name"`             // Original filename
	PassengerName   string         `json:"passenger_name"`                        // Nama penumpang dari receipt
	RouteOrLocation string         `json:"route_or_location"`                     // SUB-HLP or Hotel name
	ParsedData      string         `gorm:"type:jsonb" json:"parsed_data,omitempty"` // Data hasil parsing OCR (JSON)
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
