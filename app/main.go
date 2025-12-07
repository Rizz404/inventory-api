// Package main provides the main entry point for the Inventory Management API.
//
//	@title			Inventory Management API
//	@version		1.0
//	@description	A comprehensive inventory management API with JWT authentication, multi-language support, and CRUD operations for assets, users, and locations.
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:5000
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//
//	@tag.name			Authentication
//	@tag.description	Authentication related endpoints
//	@tag.name			Users
//	@tag.description	User management endpoints
//	@tag.name			Categories
//	@tag.description	Category management endpoints
//	@tag.name			Locations
//	@tag.description	Location management endpoints
//	@tag.name			Assets
//	@tag.description	Asset management endpoints
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Rizz404/inventory-api/config"
	_ "github.com/Rizz404/inventory-api/docs"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/rest"
	"github.com/Rizz404/inventory-api/services/asset"
	assetMovement "github.com/Rizz404/inventory-api/services/asset_movement"
	"github.com/Rizz404/inventory-api/services/auth"
	"github.com/Rizz404/inventory-api/services/category"
	issueReport "github.com/Rizz404/inventory-api/services/issue_report"
	"github.com/Rizz404/inventory-api/services/location"
	maintenanceRecord "github.com/Rizz404/inventory-api/services/maintenance_record"
	maintenanceSchedule "github.com/Rizz404/inventory-api/services/maintenance_schedule"
	"github.com/Rizz404/inventory-api/services/notification"
	scanLog "github.com/Rizz404/inventory-api/services/scan_log"
	"github.com/Rizz404/inventory-api/services/user"
	"github.com/common-nighthawk/go-figure"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}
}

func main() {
	// *===================================ENV===================================*
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":5000"
		log.Printf("ADDR environment variable not set, using default :5000")
	}

	// *===================================DATABASE===================================*
	db := config.InitializeDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	// *===================================EXTERNAL CLIENTS===================================*
	clients := config.InitializeClients()

	// *===================================REPOSITORY===================================*
	userRepository := postgresql.NewUserRepository(db)
	categoryRepository := postgresql.NewCategoryRepository(db)
	locationRepository := postgresql.NewLocationRepository(db)
	assetRepository := postgresql.NewAssetRepository(db)
	scanLogRepository := postgresql.NewScanLogRepository(db)
	notificationRepository := postgresql.NewNotificationRepository(db)
	issueReportRepository := postgresql.NewIssueReportRepository(db)
	assetMovementRepository := postgresql.NewAssetMovementRepository(db)
	maintenanceScheduleRepository := postgresql.NewMaintenanceScheduleRepository(db)
	maintenanceRecordRepository := postgresql.NewMaintenanceRecordRepository(db)

	// *===================================SERVICE===================================*
	authService := auth.NewService(userRepository)
	userService := user.NewService(userRepository, clients.Cloudinary)
	notificationService := notification.NewService(notificationRepository, userRepository, clients.FCM)
	categoryService := category.NewService(categoryRepository, notificationService, userRepository)
	locationService := location.NewService(locationRepository, notificationService, userRepository)
	assetService := asset.NewService(assetRepository, clients.Cloudinary, notificationService, categoryService, userRepository)
	scanLogService := scanLog.NewService(scanLogRepository)
	issueReportService := issueReport.NewService(issueReportRepository, notificationService, assetService, userRepository)
	assetMovementService := assetMovement.NewService(assetMovementRepository, assetService, locationService, userService, notificationService)
	maintenanceScheduleService := maintenanceSchedule.NewService(maintenanceScheduleRepository, assetService, userService, notificationService)
	maintenanceRecordService := maintenanceRecord.NewService(maintenanceRecordRepository, assetService, userService, notificationService)

	// *===================================CRON SERVICE===================================*
	assetCronService := asset.NewCronService(assetRepository, notificationService)
	if err := assetCronService.Start(); err != nil {
		log.Fatalf("Failed to start asset cron service: %v", err)
	}
	defer assetCronService.Stop()

	// *===================================SERVER CONFIG===================================*
	app := fiber.New(fiber.Config{
		AppName:       "Project Management Api",
		BodyLimit:     10 * 1024 * 1024,
		CaseSensitive: true,
		// StrictRouting: true, // ! berbahaya asw
	})

	// *===================================MIDDLEWARE===================================*
	app.Use(recovermw.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
		// AllowCredentials: true,
	}))
	app.Use(helmet.New())
	app.Use(favicon.New())
	app.Use(logger.New())
	app.Use(healthcheck.New())

	// *===================================ROUTES===================================*
	app.Get("/metrics", monitor.New())

	app.Get("/", func(c *fiber.Ctx) error {
		banner := figure.NewFigure("Inventory API", "", true).String()
		message := "Welcome to the Inventory Management API.\nDocs: /docs/index.html\nUse /api/v1/* endpoints such as /api/v1/auth/login, /api/v1/users, /api/v1/assets."
		html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Inventory API</title>
	<style>
		body {
			background-color: #0d1117;
			color: #c9d1d9;
			font-family: 'Courier New', monospace;
			padding: 20px;
			margin: 0;
		}
		pre {
			color: #58a6ff;
			font-size: 14px;
			line-height: 1.5;
			white-space: pre;
			margin: 0;
		}
		.message {
			color: #8b949e;
			margin-top: 20px;
			line-height: 1.8;
		}
		.message a {
			color: #58a6ff;
			text-decoration: none;
		}
		.message a:hover {
			text-decoration: underline;
		}
	</style>
</head>
<body>
	<pre>` + banner + `</pre>
	<div class="message">` + message + `</div>
</body>
</html>`
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "You are at /api. Jump into versioned routes under /api/v1/*.",
			"examples": []string{"/api/v1/auth/login", "/api/v1/users", "/api/v1/assets", "/api/v1/categories"},
			"docs":     "/docs/index.html",
		})
	})

	// Swagger documentation route
	app.Get("/docs/*", swagger.New(swagger.Config{}))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":   "You are at /api/v1. Use the resource endpoints below.",
			"resources": []string{"/api/v1/auth/login", "/api/v1/users", "/api/v1/categories", "/api/v1/locations", "/api/v1/assets", "/api/v1/notifications", "/api/v1/issue-reports", "/api/v1/asset-movements", "/api/v1/maintenance-schedules", "/api/v1/maintenance-records", "/api/v1/scan-logs"},
			"docs":      "/docs/index.html",
		})
	})

	rest.NewAuthHandler(v1, authService)
	rest.NewUserHandler(v1, userService)
	rest.NewCategoryHandler(v1, categoryService)
	rest.NewLocationHandler(v1, locationService)
	rest.NewAssetHandler(v1, assetService)
	rest.NewScanLogHandler(v1, scanLogService)
	rest.NewNotificationHandler(v1, notificationService)
	rest.NewIssueReportHandler(v1, issueReportService)
	rest.NewAssetMovementHandler(v1, assetMovementService)
	rest.NewMaintenanceScheduleHandler(v1, maintenanceScheduleService)
	rest.NewMaintenanceRecordHandler(v1, maintenanceRecordService)

	// *===================================SERVER===================================*
	log.Printf("server running on http://localhost%s", addr)

	if err := app.Listen(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
