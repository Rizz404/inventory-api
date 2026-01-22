package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitializeDatabase initializes database connection
func InitializeDatabase() *gorm.DB {
	DSN := os.Getenv("DSN")
	if DSN == "" {
		log.Fatalf("DSN environment variable not set")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * 3 / 2, // 1.5 seconds
			IgnoreRecordNotFoundError: true,
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("failed to open connection to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic database object: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Printf("successfully connected to database")

	return db
}
