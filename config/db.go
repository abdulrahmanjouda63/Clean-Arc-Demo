package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes and returns a *gorm.DB based on config (mysql by default)
func InitDB(cfg *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	
	switch cfg.DB.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect mysql: %w", err)
		}
	default:
		// default to mysql if driver is not specified or recognized
		db, err = gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect default mysql: %w", err)
		}
	}
	
	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	return db, nil
}
