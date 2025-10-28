package database

import (
	"fmt"
	"log"
	"perjalanan-dinas/backend/config"
	"perjalanan-dinas/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.Admin{},
		&models.Position{},
		&models.Employee{},
		&models.TravelRequest{},
		&models.TravelRequestEmployee{},
		&models.TravelReport{},
		&models.VisitProof{},
		&models.NumberingConfig{},
		&models.RepresentativeConfig{},
		&models.AtCostClaim{},
		&models.AtCostClaimItem{},
		&models.AtCostReceipt{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")

	// Seed positions if not exists
	if err := seedPositions(); err != nil {
		log.Printf("Warning: failed to seed positions: %v", err)
	}

	// Create default admin if not exists
	if err := createDefaultAdmin(cfg); err != nil {
		log.Printf("Warning: failed to create default admin: %v", err)
	}

	// Initialize numbering config if not exists
	if err := initializeNumberingConfig(); err != nil {
		log.Printf("Warning: failed to initialize numbering config: %v", err)
	}

	// Initialize representative config if not exists
	if err := initializeRepresentativeConfig(); err != nil {
		log.Printf("Warning: failed to initialize representative config: %v", err)
	}

	return nil
}

func createDefaultAdmin(cfg *config.Config) error {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := models.Admin{
			Username: cfg.AdminUsername,
			Password: string(hashedPassword),
		}

		if err := DB.Create(&admin).Error; err != nil {
			return err
		}

		log.Printf("Default admin created: %s", cfg.AdminUsername)
	}

	return nil
}

func seedPositions() error {
	var count int64
	DB.Model(&models.Position{}).Count(&count)

	if count == 0 {
		for _, posData := range AllPositions {
			position := models.Position{
				Title:                    posData.Title,
				Code:                     posData.Code,
				Level:                    posData.Level,
				AllowanceInProvince:      posData.AllowanceInProvince,
				AllowanceOutsideProvince: posData.AllowanceOutsideProvince,
				AllowanceAbroad:          posData.AllowanceAbroad,
			}
			if err := DB.Create(&position).Error; err != nil {
				return err
			}
		}
		log.Printf("Positions seeded successfully: %d positions", len(AllPositions))
	}

	return nil
}

func initializeNumberingConfig() error {
	var count int64
	DB.Model(&models.NumberingConfig{}).Count(&count)

	if count == 0 {
		config := models.NumberingConfig{
			LastRequestSequence: 0,
			LastReportSequence:  0,
		}

		if err := DB.Create(&config).Error; err != nil {
			return err
		}

		log.Println("Numbering config initialized")
	}

	return nil
}

func initializeRepresentativeConfig() error {
	var count int64
	DB.Model(&models.RepresentativeConfig{}).Count(&count)

	if count == 0 {
		config := models.RepresentativeConfig{
			Name:     "M. MACHFUD HIDAYAT",
			Position: "Vice President",
			IsActive: true,
		}

		if err := DB.Create(&config).Error; err != nil {
			return err
		}

		log.Println("Representative config initialized with default values")
	}

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
