package config

import (
	"fmt"
	"log"

	"gadget-marketplace/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *Config) *gorm.DB {
	var db *gorm.DB
	var err error

	if cfg.DBDriver == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Warning: Failed to connect to PostgreSQL (%v). Falling back to SQLite local database...", err)
			db, err = gorm.Open(sqlite.Open("gadget_marketplace.db"), &gorm.Config{})
			if err != nil {
				log.Fatalf("Fatal: Failed to connect to fallback SQLite database: %v", err)
			}
			log.Println("Connected successfully to SQLite local database (gadget_marketplace.db)")
		} else {
			log.Println("Connected successfully to PostgreSQL database!")
		}
	} else {
		db, err = gorm.Open(sqlite.Open("gadget_marketplace.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Fatal: Failed to connect to SQLite database: %v", err)
		}
		log.Println("Connected successfully to SQLite database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	if err != nil {
		log.Fatalf("Fatal: Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Database Auto-Migration completed successfully!")
	return db
}
