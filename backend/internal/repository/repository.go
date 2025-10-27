package repository

import (
	"errors"
	"fmt"
	"perjalanan-dinas/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Admin operations
func (r *Repository) GetAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	return &admin, err
}

// Position operations
func (r *Repository) GetAllPositions() ([]models.Position, error) {
	var positions []models.Position
	err := r.db.Find(&positions).Error
	return positions, err
}

func (r *Repository) GetPositionByID(id uint) (*models.Position, error) {
	var position models.Position
	err := r.db.First(&position, id).Error
	return &position, err
}

// Employee operations
func (r *Repository) CreateEmployee(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *Repository) GetAllEmployees() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Preload("Position").Find(&employees).Error
	return employees, err
}

func (r *Repository) GetEmployeeByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Preload("Position").First(&employee, id).Error
	return &employee, err
}

func (r *Repository) UpdateEmployee(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *Repository) DeleteEmployee(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

// NumberingConfig operations
func (r *Repository) GetNumberingConfig() (*models.NumberingConfig, error) {
	var config models.NumberingConfig
	err := r.db.First(&config).Error
	return &config, err
}

func (r *Repository) UpdateNumberingConfig(config *models.NumberingConfig) error {
	return r.db.Save(config).Error
}

func (r *Repository) GetNextRequestSequence() (int, error) {
	config, err := r.GetNumberingConfig()
	if err != nil {
		return 0, err
	}

	config.LastRequestSequence++
	if err := r.UpdateNumberingConfig(config); err != nil {
		return 0, err
	}

	return config.LastRequestSequence, nil
}

func (r *Repository) GetNextReportSequence() (int, error) {
	config, err := r.GetNumberingConfig()
	if err != nil {
		return 0, err
	}

	config.LastReportSequence++
	if err := r.UpdateNumberingConfig(config); err != nil {
		return 0, err
	}

	return config.LastReportSequence, nil
}

// TravelRequest operations
func (r *Repository) CreateTravelRequest(request *models.TravelRequest) error {
	return r.db.Create(request).Error
}

func (r *Repository) GetAllTravelRequests() ([]models.TravelRequest, error) {
	var requests []models.TravelRequest
	err := r.db.Preload("TravelRequestEmployees.Employee.Position").Find(&requests).Error
	return requests, err
}

func (r *Repository) GetTravelRequestByID(id uint) (*models.TravelRequest, error) {
	var request models.TravelRequest
	err := r.db.Preload("TravelRequestEmployees.Employee.Position").
		Preload("TravelReport").
		First(&request, id).Error
	return &request, err
}

func (r *Repository) UpdateTravelRequest(request *models.TravelRequest) error {
	return r.db.Save(request).Error
}

func (r *Repository) DeleteTravelRequest(id uint) error {
	return r.db.Delete(&models.TravelRequest{}, id).Error
}

// TravelReport operations
func (r *Repository) CreateTravelReport(report *models.TravelReport) error {
	return r.db.Create(report).Error
}

func (r *Repository) GetTravelReportByRequestID(requestID uint) (*models.TravelReport, error) {
	var report models.TravelReport
	err := r.db.Preload("VisitProofs").Where("travel_request_id = ?", requestID).First(&report).Error
	return &report, err
}

func (r *Repository) UpdateTravelReport(report *models.TravelReport) error {
	return r.db.Save(report).Error
}

// VisitProof operations
func (r *Repository) CreateVisitProof(proof *models.VisitProof) error {
	return r.db.Create(proof).Error
}

func (r *Repository) GetVisitProofsByReportID(reportID uint) ([]models.VisitProof, error) {
	var proofs []models.VisitProof
	err := r.db.Where("travel_report_id = ?", reportID).Find(&proofs).Error
	return proofs, err
}

func (r *Repository) UpdateVisitProof(proof *models.VisitProof) error {
	return r.db.Save(proof).Error
}

func (r *Repository) DeleteVisitProof(id uint) error {
	return r.db.Delete(&models.VisitProof{}, id).Error
}

// Calculate duration days
func CalculateDurationDays(departure, returnDate time.Time) int {
	duration := returnDate.Sub(departure)
	days := int(duration.Hours() / 24)
	if days <= 0 {
		return 1
	}
	return days + 1 // Include both departure and return day
}

// Generate request number: 064/{seq}/DIB/{code}/NOTA
func (r *Repository) GenerateRequestNumber(employeeID uint) (string, error) {
	employee, err := r.GetEmployeeByID(employeeID)
	if err != nil {
		return "", err
	}

	position, err := r.GetPositionByID(employee.PositionID)
	if err != nil {
		return "", errors.New("position not found for employee")
	}

	seq, err := r.GetNextRequestSequence()
	if err != nil {
		return "", err
	}

	// Format: 064/{seq}/DIB/{code}/NOTA
	requestNumber := formatRequestNumber(seq, position.Code)
	return requestNumber, nil
}

// Generate report number: 064/ /DIB/{code}/NOTA
func (r *Repository) GenerateReportNumber(employeeID uint) (string, error) {
	employee, err := r.GetEmployeeByID(employeeID)
	if err != nil {
		return "", err
	}

	position, err := r.GetPositionByID(employee.PositionID)
	if err != nil {
		return "", errors.New("position not found for employee")
	}

	// Format: 064/ /DIB/{code}/NOTA (with space instead of sequence)
	reportNumber := formatReportNumber(position.Code)
	return reportNumber, nil
}

func formatRequestNumber(seq int, code string) string {
	return "064/" + padLeft(seq, 4) + "/DIB/" + code + "/NOTA"
}

func formatReportNumber(code string) string {
	return "064/    /DIB/" + code + "/NOTA"
}

func padLeft(num int, length int) string {
	numStr := fmt.Sprintf("%d", num)
	if len(numStr) >= length {
		return numStr
	}
	return fmt.Sprintf("%0"+fmt.Sprintf("%d", length)+"d", num)
}

// GetDB returns the database instance (for handler access)
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

// CreateTravelRequestEmployee creates a new travel request employee relation
func (r *Repository) CreateTravelRequestEmployee(empRel *models.TravelRequestEmployee) error {
	return r.db.Create(empRel).Error
}

// Representative Config operations
func (r *Repository) GetActiveRepresentativeConfig() (*models.RepresentativeConfig, error) {
	var config models.RepresentativeConfig
	if err := r.db.Where("is_active = ?", true).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *Repository) GetRepresentativeConfig() (*models.RepresentativeConfig, error) {
	var config models.RepresentativeConfig
	if err := r.db.First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *Repository) UpdateRepresentativeConfig(config *models.RepresentativeConfig) error {
	return r.db.Save(config).Error
}
