package database

// Position represents a job position with its code and allowance rates
type PositionData struct {
	Title                    string
	Code                     string
	Level                    string // Jr. Officer, Officer, Senior Officer, AVP
	AllowanceInProvince      int
	AllowanceOutsideProvince int
	AllowanceAbroad          int
}

// AllPositions contains all available positions with their allowance rates
var AllPositions = []PositionData{
	// DPEB - Pengembangan Ekosistem & Bisnis Digital
	{
		Title:                    "Jr. Officer Pengembangan Bisnis Digital",
		Code:                     "DPEB",
		Level:                    "Jr. Officer",
		AllowanceInProvince:      10000,
		AllowanceOutsideProvince: 20000,
		AllowanceAbroad:          30000,
	},
	{
		Title:                    "Officer Pengembangan Bisnis Digital",
		Code:                     "DPEB",
		Level:                    "Officer",
		AllowanceInProvince:      15000,
		AllowanceOutsideProvince: 30000,
		AllowanceAbroad:          45000,
	},
	{
		Title:                    "Senior Officer Pengembangan Bisnis Digital",
		Code:                     "DPEB",
		Level:                    "Senior Officer",
		AllowanceInProvince:      20000,
		AllowanceOutsideProvince: 40000,
		AllowanceAbroad:          60000,
	},
	{
		Title:                    "Jr. Officer Pengembangan Ekosistem Digital & Kerjasama Bisnis",
		Code:                     "DPEB",
		Level:                    "Jr. Officer",
		AllowanceInProvince:      10000,
		AllowanceOutsideProvince: 20000,
		AllowanceAbroad:          30000,
	},
	{
		Title:                    "Officer Pengembangan Ekosistem Digital & Kerjasama Bisnis",
		Code:                     "DPEB",
		Level:                    "Officer",
		AllowanceInProvince:      15000,
		AllowanceOutsideProvince: 30000,
		AllowanceAbroad:          45000,
	},
	{
		Title:                    "Senior Officer Pengembangan Ekosistem Digital & Kerjasama Bisnis",
		Code:                     "DPEB",
		Level:                    "Senior Officer",
		AllowanceInProvince:      20000,
		AllowanceOutsideProvince: 40000,
		AllowanceAbroad:          60000,
	},
	{
		Title:                    "AVP Pengembangan Ekosistem & Bisnis Digital",
		Code:                     "DPEB",
		Level:                    "AVP",
		AllowanceInProvince:      25000,
		AllowanceOutsideProvince: 50000,
		AllowanceAbroad:          75000,
	},

	// DPDB - Digital Banking & e-Channel
	{
		Title:                    "Jr. Officer Pengembangan Digital Banking & e-Channel",
		Code:                     "DPDB",
		Level:                    "Jr. Officer",
		AllowanceInProvince:      10000,
		AllowanceOutsideProvince: 20000,
		AllowanceAbroad:          30000,
	},
	{
		Title:                    "Officer Pengembangan Digital Banking & e-Channel",
		Code:                     "DPDB",
		Level:                    "Officer",
		AllowanceInProvince:      15000,
		AllowanceOutsideProvince: 30000,
		AllowanceAbroad:          45000,
	},
	{
		Title:                    "Senior Officer Pengembangan Digital Banking & e-Channel",
		Code:                     "DPDB",
		Level:                    "Senior Officer",
		AllowanceInProvince:      20000,
		AllowanceOutsideProvince: 40000,
		AllowanceAbroad:          60000,
	},
	{
		Title:                    "Jr. Officer Manajemen Aplikasi",
		Code:                     "DPDB",
		Level:                    "Jr. Officer",
		AllowanceInProvince:      10000,
		AllowanceOutsideProvince: 20000,
		AllowanceAbroad:          30000,
	},
	{
		Title:                    "Officer Manajemen Aplikasi",
		Code:                     "DPDB",
		Level:                    "Officer",
		AllowanceInProvince:      15000,
		AllowanceOutsideProvince: 30000,
		AllowanceAbroad:          45000,
	},
	{
		Title:                    "Senior Officer Manajemen Aplikasi",
		Code:                     "DPDB",
		Level:                    "Senior Officer",
		AllowanceInProvince:      20000,
		AllowanceOutsideProvince: 40000,
		AllowanceAbroad:          60000,
	},
	{
		Title:                    "AVP Manajemen Aplikasi Digital Banking",
		Code:                     "DPDB",
		Level:                    "AVP",
		AllowanceInProvince:      25000,
		AllowanceOutsideProvince: 50000,
		AllowanceAbroad:          75000,
	},

	// DDBE - QA & Monitoring Digital Banking
	{
		Title:                    "Jr. Officer IT & Digital Banking Quality Assurance",
		Code:                     "DDBE",
		Level:                    "Jr. Officer",
		AllowanceInProvince:      10000,
		AllowanceOutsideProvince: 20000,
		AllowanceAbroad:          30000,
	},
	{
		Title:                    "Officer IT & Digital Banking Quality Assurance",
		Code:                     "DDBE",
		Level:                    "Officer",
		AllowanceInProvince:      15000,
		AllowanceOutsideProvince: 30000,
		AllowanceAbroad:          45000,
	},
	{
		Title:                    "Senior Officer IT & Digital Banking Quality Assurance",
		Code:                     "DDBE",
		Level:                    "Senior Officer",
		AllowanceInProvince:      20000,
		AllowanceOutsideProvince: 40000,
		AllowanceAbroad:          60000,
	},
	{
		Title:                    "Jr. Officer Monitoring & Evaluasi Digital Banking",
		Code:                     "DDBE",
		Level:                    "Jr. Officer",
		AllowanceInProvince:      10000,
		AllowanceOutsideProvince: 20000,
		AllowanceAbroad:          30000,
	},
	{
		Title:                    "Officer Monitoring & Evaluasi Digital Banking",
		Code:                     "DDBE",
		Level:                    "Officer",
		AllowanceInProvince:      15000,
		AllowanceOutsideProvince: 30000,
		AllowanceAbroad:          45000,
	},
	{
		Title:                    "Senior Officer Monitoring & Evaluasi Digital Banking",
		Code:                     "DDBE",
		Level:                    "Senior Officer",
		AllowanceInProvince:      20000,
		AllowanceOutsideProvince: 40000,
		AllowanceAbroad:          60000,
	},
	{
		Title:                    "AVP QA & Monitoring Digital Banking",
		Code:                     "DDBE",
		Level:                    "AVP",
		AllowanceInProvince:      25000,
		AllowanceOutsideProvince: 50000,
		AllowanceAbroad:          75000,
	},
}
