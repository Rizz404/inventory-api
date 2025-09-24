package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}
}

func main() {
	// *===================================ENV===================================*
	DSN := os.Getenv("DSN")
	if DSN == "" {
		log.Fatalf("DSN environment variable not set")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":5000"
		log.Printf("ADDR environment variable not set, using default :5000")
	}

	// * FCM Configuration
	// enableFCM := os.Getenv("ENABLE_FCM") == "true"

	// *===================================DATABASE===================================*
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("failed to open connection to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Printf("successfully connected to database")

	// *===================================CLOUDINARY CLIENT===================================*
	var cloudinaryClient *cloudinary.Client
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL != "" {
		var err error
		cloudinaryClient, err = cloudinary.NewClientFromURL(cloudinaryURL)
		if err != nil {
			log.Printf("Warning: Failed to initialize Cloudinary client: %v. File upload will be disabled.", err)
		} else {
			log.Printf("Cloudinary client initialized successfully")
		}
	} else {
		// Try individual environment variables
		cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
		apiKey := os.Getenv("CLOUDINARY_API_KEY")
		apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

		if cloudName != "" && apiKey != "" && apiSecret != "" {
			var err error
			cloudinaryClient, err = cloudinary.NewClient(cloudName, apiKey, apiSecret)
			if err != nil {
				log.Printf("Warning: Failed to initialize Cloudinary client: %v. File upload will be disabled.", err)
			} else {
				log.Printf("Cloudinary client initialized successfully")
			}
		} else {
			log.Printf("Cloudinary credentials not provided. File upload will be disabled.")
		}
	}

	// *===================================FCM CLIENT===================================*
	// var fcmClient *fcm.Client
	// if enableFCM {
	// 	credentialsMap := map[string]string{
	// 		"type":                        os.Getenv("FIREBASE_TYPE"),
	// 		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
	// 		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
	// 		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
	// 		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
	// 		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
	// 		"auth_uri":                    os.Getenv("FIREBASE_AUTH_URI"),
	// 		"token_uri":                   os.Getenv("FIREBASE_TOKEN_URI"),
	// 		"auth_provider_x509_cert_url": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
	// 		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	// 		"universe_domain":             os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
	// 	}

	// 	credentialsJSON, err := json.Marshal(credentialsMap)
	// 	if err != nil {
	// 		log.Printf("Warning: Failed to marshal Firebase credentials: %v. Firebase services will be disabled.", err)
	// 	} else {
	// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// 		defer cancel()

	// 		opt := option.WithCredentialsJSON(credentialsJSON)

	// 		app, err := firebase.NewApp(ctx, nil, opt)
	// 		if err != nil {
	// 			log.Printf("Warning: Failed to initialize Firebase app: %v. Firebase services will be disabled.", err)
	// 		} else {
	// 			// * Inisialisasi FCM Client
	// 			client, err := app.Messaging(ctx)
	// 			if err != nil {
	// 				log.Printf("Warning: Failed to get FCM messaging client: %v. FCM will be disabled.", err)
	// 			} else {
	// 				fcmClient = fcm.NewClientFromMessaging(client)
	// 				log.Printf("FCM client initialized successfully")
	// 			}

	// 			// Todo: Nanti inisialisasi service dari firebase lain disini
	// 		}
	// 	}
	// } else {
	// 	log.Printf("Firebase services disabled via ENABLE_FCM environment variable")
	// }

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
	userService := user.NewService(userRepository, cloudinaryClient)
	categoryService := category.NewService(categoryRepository)
	locationService := location.NewService(locationRepository)
	assetService := asset.NewService(assetRepository, cloudinaryClient)
	scanLogService := scanLog.NewService(scanLogRepository)
	notificationService := notification.NewService(notificationRepository)
	issueReportService := issueReport.NewService(issueReportRepository)
	assetMovementService := assetMovement.NewService(assetMovementRepository, assetService, locationService, userService)
	maintenanceScheduleService := maintenanceSchedule.NewService(maintenanceScheduleRepository, assetService, userService)
	maintenanceRecordService := maintenanceRecord.NewService(maintenanceRecordRepository, assetService, userService)

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

	api := app.Group("/api")
	v1 := api.Group("/v1")

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
