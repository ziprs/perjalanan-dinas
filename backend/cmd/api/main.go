package main

import (
	"log"
	"perjalanan-dinas/backend/config"
	"perjalanan-dinas/backend/internal/database"
	"perjalanan-dinas/backend/internal/handlers"
	"perjalanan-dinas/backend/internal/middleware"
	"perjalanan-dinas/backend/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize JWT secret
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Initialize repository
	repo := repository.NewRepository(database.GetDB())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(repo)
	employeeHandler := handlers.NewEmployeeHandler(repo)
	positionHandler := handlers.NewPositionHandler(repo)
	cityHandler := handlers.NewCityHandler()
	travelRequestHandler := handlers.NewTravelRequestHandler(repo)
	travelReportHandler := handlers.NewTravelReportHandler(repo)
	pdfHandler := handlers.NewPDFHandler(repo)
	excelHandler := handlers.NewExcelHandler(repo)
	representativeHandler := handlers.NewRepresentativeHandler(repo)
	atCostHandler := handlers.NewAtCostHandler(repo)
	healthHandler := handlers.NewHealthHandler()

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint (no /api prefix)
	router.GET("/health", healthHandler.Health)

	// Public routes
	public := router.Group("/api")
	{
		// Auth
		public.POST("/auth/login", authHandler.Login)

		// Public access for employees to create travel request
		public.GET("/employees", employeeHandler.GetAllEmployees)
		public.GET("/employees/:id", employeeHandler.GetEmployeeByID)
		public.GET("/positions", positionHandler.GetAllPositions)
		public.GET("/cities", cityHandler.GetAllCities)

		// Travel requests - public for employees to submit
		public.POST("/travel-requests", travelRequestHandler.CreateTravelRequest)
		public.GET("/travel-requests/stats/employees", travelRequestHandler.GetEmployeeSPDStats)
		public.GET("/travel-requests/:id", travelRequestHandler.GetTravelRequestByID)
		public.GET("/travel-requests", travelRequestHandler.GetAllTravelRequests)

		// At-Cost claims - public for employees to submit
		public.POST("/at-cost/upload-receipt", atCostHandler.UploadReceipt)
		public.POST("/at-cost/claims", atCostHandler.CreateAtCostClaim)
		public.GET("/at-cost/claims", atCostHandler.GetAllAtCostClaims)
		public.GET("/at-cost/claims/:id", atCostHandler.GetAtCostClaim)
		public.GET("/at-cost/receipts/:receipt_id/download", atCostHandler.DownloadReceipt)

		// PDF downloads - public
		public.GET("/pdf/nota-permintaan/:id", pdfHandler.DownloadNotaPermintaan)
		public.GET("/pdf/berita-acara/:id", pdfHandler.DownloadBeritaAcara)
		public.GET("/pdf/combined/:id", pdfHandler.DownloadCombinedPDF)
		public.GET("/pdf/nota-atcost/:id", atCostHandler.DownloadNotaAtCost)
		public.GET("/pdf/combined-atcost/:id", atCostHandler.DownloadCombinedAtCost)
	}

	// Protected routes (require authentication)
	protected := router.Group("/api/admin")
	protected.Use(middleware.AuthMiddleware())
	{
		// Employee management
		protected.POST("/employees", employeeHandler.CreateEmployee)
		protected.PUT("/employees/:id", employeeHandler.UpdateEmployee)
		protected.DELETE("/employees/:id", employeeHandler.DeleteEmployee)

		// Travel requests management
		protected.GET("/travel-requests", travelRequestHandler.GetAllTravelRequests)
		protected.DELETE("/travel-requests/:id", travelRequestHandler.DeleteTravelRequest)

		// Travel reports
		protected.POST("/travel-reports", travelReportHandler.CreateTravelReport)
		protected.GET("/travel-reports/:request_id", travelReportHandler.GetTravelReportByRequestID)

		// Excel export
		protected.GET("/excel/monthly-allowance", excelHandler.ExportMonthlyAllowance)

		// Representative config
		protected.GET("/representative-config", representativeHandler.GetRepresentativeConfig)
		protected.PUT("/representative-config", representativeHandler.UpdateRepresentativeConfig)

		// At-Cost claims
		protected.POST("/at-cost/upload-receipt", atCostHandler.UploadReceipt)
		protected.POST("/at-cost/claims", atCostHandler.CreateAtCostClaim)
		protected.GET("/at-cost/claims", atCostHandler.GetAllAtCostClaims)
		protected.GET("/at-cost/claims/:id", atCostHandler.GetAtCostClaim)
		protected.GET("/at-cost/claims/travel-request/:travel_request_id", atCostHandler.GetAtCostClaimByTravelRequest)
		protected.PUT("/at-cost/claims/:id/status", atCostHandler.UpdateClaimStatus)
		protected.DELETE("/at-cost/claims/:id", atCostHandler.DeleteAtCostClaim)
		protected.GET("/at-cost/receipts/:receipt_id/download", atCostHandler.DownloadReceipt)
		protected.POST("/at-cost/parse-manual", atCostHandler.ParseReceiptManual)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
