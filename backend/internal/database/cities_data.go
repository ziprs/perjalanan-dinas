package database

// CityData represents a city or country destination
type CityData struct {
	Name            string
	DestinationType string // "in_province", "outside_province", "abroad"
}

// AllCities contains Indonesian cities grouped by destination type and foreign countries
var AllCities = []CityData{
	// Dalam Provinsi Jawa Timur
	{"Surabaya", "in_province"},
	{"Malang", "in_province"},
	{"Sidoarjo", "in_province"},
	{"Gresik", "in_province"},
	{"Mojokerto", "in_province"},
	{"Pasuruan", "in_province"},
	{"Probolinggo", "in_province"},
	{"Blitar", "in_province"},
	{"Kediri", "in_province"},
	{"Madiun", "in_province"},
	{"Banyuwangi", "in_province"},
	{"Jember", "in_province"},
	{"Situbondo", "in_province"},
	{"Bondowoso", "in_province"},
	{"Lumajang", "in_province"},
	{"Tulungagung", "in_province"},
	{"Nganjuk", "in_province"},
	{"Jombang", "in_province"},
	{"Bojonegoro", "in_province"},
	{"Tuban", "in_province"},
	{"Lamongan", "in_province"},
	{"Bangkalan", "in_province"},
	{"Sampang", "in_province"},
	{"Pamekasan", "in_province"},
	{"Sumenep", "in_province"},
	{"Ngawi", "in_province"},
	{"Magetan", "in_province"},
	{"Ponorogo", "in_province"},
	{"Pacitan", "in_province"},
	{"Trenggalek", "in_province"},
	{"Batu", "in_province"},

	// Luar Provinsi - Pulau Jawa
	{"Jakarta", "outside_province"},
	{"Bandung", "outside_province"},
	{"Semarang", "outside_province"},
	{"Yogyakarta", "outside_province"},
	{"Solo (Surakarta)", "outside_province"},
	{"Bekasi", "outside_province"},
	{"Tangerang", "outside_province"},
	{"Depok", "outside_province"},
	{"Bogor", "outside_province"},
	{"Cirebon", "outside_province"},
	{"Sukabumi", "outside_province"},
	{"Tasikmalaya", "outside_province"},
	{"Purwokerto", "outside_province"},
	{"Tegal", "outside_province"},
	{"Pekalongan", "outside_province"},
	{"Magelang", "outside_province"},
	{"Salatiga", "outside_province"},
	{"Serang", "outside_province"},
	{"Cilegon", "outside_province"},

	// Luar Provinsi - Sumatera
	{"Medan", "outside_province"},
	{"Palembang", "outside_province"},
	{"Pekanbaru", "outside_province"},
	{"Padang", "outside_province"},
	{"Bandar Lampung", "outside_province"},
	{"Jambi", "outside_province"},
	{"Bengkulu", "outside_province"},
	{"Banda Aceh", "outside_province"},
	{"Batam", "outside_province"},
	{"Dumai", "outside_province"},
	{"Bukittinggi", "outside_province"},
	{"Tanjung Pinang", "outside_province"},

	// Luar Provinsi - Kalimantan
	{"Balikpapan", "outside_province"},
	{"Samarinda", "outside_province"},
	{"Banjarmasin", "outside_province"},
	{"Pontianak", "outside_province"},
	{"Palangkaraya", "outside_province"},
	{"Tarakan", "outside_province"},
	{"Bontang", "outside_province"},
	{"Singkawang", "outside_province"},

	// Luar Provinsi - Sulawesi
	{"Makassar", "outside_province"},
	{"Manado", "outside_province"},
	{"Palu", "outside_province"},
	{"Kendari", "outside_province"},
	{"Gorontalo", "outside_province"},
	{"Mamuju", "outside_province"},

	// Luar Provinsi - Bali & Nusa Tenggara
	{"Denpasar", "outside_province"},
	{"Mataram", "outside_province"},
	{"Kupang", "outside_province"},

	// Luar Provinsi - Maluku & Papua
	{"Ambon", "outside_province"},
	{"Ternate", "outside_province"},
	{"Jayapura", "outside_province"},
	{"Sorong", "outside_province"},
	{"Manokwari", "outside_province"},

	// Luar Negeri - Asia Tenggara
	{"Singapura", "abroad"},
	{"Kuala Lumpur, Malaysia", "abroad"},
	{"Bangkok, Thailand", "abroad"},
	{"Manila, Filipina", "abroad"},
	{"Hanoi, Vietnam", "abroad"},
	{"Ho Chi Minh, Vietnam", "abroad"},
	{"Yangon, Myanmar", "abroad"},
	{"Phnom Penh, Kamboja", "abroad"},
	{"Vientiane, Laos", "abroad"},
	{"Bandar Seri Begawan, Brunei", "abroad"},

	// Luar Negeri - Asia Timur
	{"Tokyo, Jepang", "abroad"},
	{"Seoul, Korea Selatan", "abroad"},
	{"Beijing, China", "abroad"},
	{"Shanghai, China", "abroad"},
	{"Hong Kong", "abroad"},
	{"Taipei, Taiwan", "abroad"},

	// Luar Negeri - Asia Selatan & Tengah
	{"New Delhi, India", "abroad"},
	{"Mumbai, India", "abroad"},
	{"Dhaka, Bangladesh", "abroad"},
	{"Karachi, Pakistan", "abroad"},
	{"Colombo, Sri Lanka", "abroad"},

	// Luar Negeri - Timur Tengah
	{"Dubai, UAE", "abroad"},
	{"Abu Dhabi, UAE", "abroad"},
	{"Riyadh, Arab Saudi", "abroad"},
	{"Jeddah, Arab Saudi", "abroad"},
	{"Doha, Qatar", "abroad"},
	{"Kuwait City, Kuwait", "abroad"},
	{"Muscat, Oman", "abroad"},

	// Luar Negeri - Eropa
	{"London, Inggris", "abroad"},
	{"Paris, Prancis", "abroad"},
	{"Berlin, Jerman", "abroad"},
	{"Amsterdam, Belanda", "abroad"},
	{"Brussels, Belgia", "abroad"},
	{"Roma, Italia", "abroad"},
	{"Madrid, Spanyol", "abroad"},
	{"Moskow, Rusia", "abroad"},
	{"Wina, Austria", "abroad"},
	{"Zurich, Swiss", "abroad"},

	// Luar Negeri - Amerika
	{"New York, Amerika Serikat", "abroad"},
	{"Washington DC, Amerika Serikat", "abroad"},
	{"Los Angeles, Amerika Serikat", "abroad"},
	{"San Francisco, Amerika Serikat", "abroad"},
	{"Toronto, Kanada", "abroad"},
	{"Vancouver, Kanada", "abroad"},
	{"Mexico City, Meksiko", "abroad"},
	{"Sao Paulo, Brazil", "abroad"},
	{"Buenos Aires, Argentina", "abroad"},

	// Luar Negeri - Australia & Oseania
	{"Sydney, Australia", "abroad"},
	{"Melbourne, Australia", "abroad"},
	{"Perth, Australia", "abroad"},
	{"Auckland, Selandia Baru", "abroad"},
	{"Wellington, Selandia Baru", "abroad"},

	// Luar Negeri - Afrika
	{"Kairo, Mesir", "abroad"},
	{"Johannesburg, Afrika Selatan", "abroad"},
	{"Nairobi, Kenya", "abroad"},
	{"Lagos, Nigeria", "abroad"},
}
