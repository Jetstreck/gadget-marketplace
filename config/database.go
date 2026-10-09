package config

import (
	"fmt"
	"log"

	"gadget-marketplace/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *Config) *gorm.DB {
	var db *gorm.DB
	var err error

	if cfg.DBDriver == "postgres" {
		if cfg.DBName != "postgres" {
			createPostgresDBIfNotExists(cfg)
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			PrepareStmt: false,
		})

		if err != nil {
			log.Fatalf("Fatal: Failed to connect to PostgreSQL database '%s': %v", cfg.DBName, err)
		}
		log.Println("Connected successfully to PostgreSQL database!")
	} else {
		log.Fatalf("Fatal: Invalid DB_DRIVER configuration: %s", cfg.DBDriver)
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	if err != nil {
		log.Fatalf("Fatal: Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Database Auto-Migration completed successfully!")
	return db
}

func createPostgresDBIfNotExists(cfg *Config) {
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  defaultDSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Printf("Note: Could not connect to default postgres DB for auto-creation: %v", err)
		return
	}

	var count int
	db.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", cfg.DBName).Scan(&count)
	if count == 0 {
		log.Printf("Database '%s' does not exist in PostgreSQL. Auto-creating database '%s'...", cfg.DBName, cfg.DBName)
		if err := db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\";", cfg.DBName)).Error; err != nil {
			log.Printf("Warning: Auto-creation of database failed: %v", err)
		} else {
			log.Printf("Successfully created PostgreSQL database '%s'!", cfg.DBName)
		}
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}
