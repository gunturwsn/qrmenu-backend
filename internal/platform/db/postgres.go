package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"qrmenu/internal/config"
	"qrmenu/internal/domain"
)

func Connect(c *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Optional: ping and set timezone on the session to verify connectivity.
	if sqlDB, err := gdb.DB(); err == nil {
		if c.DBMaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(c.DBMaxIdleConns)
		}
		if c.DBMaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(c.DBMaxOpenConns)
		}
		if c.DBConnMaxLifeSec > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(c.DBConnMaxLifeSec) * time.Second)
		}
		if c.DBConnMaxIdleSec > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(c.DBConnMaxIdleSec) * time.Second)
		}
		_ = sqlDB.Ping()
	}
	if err := gdb.Exec(`SET TIME ZONE 'Asia/Jakarta'`).Error; err != nil {
		log.Printf("warn: failed to SET TIME ZONE: %v", err)
	}

	return gdb
}

// Ensure required extensions exist (safe to call repeatedly).
func EnsureExtensions(db *gorm.DB) {
	_ = db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error
	// Enable uuid-ossp as well if your deployment relies on it.
	// _ = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}

func AutoMigrate(db *gorm.DB) {
	EnsureExtensions(db)

	// Order ensures foreign keys reference already-migrated entities.
	err := db.AutoMigrate(
		&domain.Tenant{},
		&domain.Table{},

		&domain.AdminUser{},

		&domain.Category{},
		&domain.Item{},

		&domain.ItemOption{},
		&domain.ItemOptionValue{},

		&domain.Order{},
		&domain.OrderItem{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Optionally set the database default timezone (applies to new sessions).
	// Safe to call multiple times.
	if derr := setDatabaseTimezone(db, "Asia/Jakarta"); derr != nil {
		log.Printf("warn: failed to set DB default timezone: %v", derr)
	}
}

func setDatabaseTimezone(db *gorm.DB, tz string) error {
	// Configure the default parameter for this database, not just the current session.
	return db.Exec(fmt.Sprintf(`ALTER DATABASE %s SET TIME ZONE '%s'`, currentDBName(db), tz)).Error
}

// currentDBName retrieves the active database name using the underlying *sql.DB.
func currentDBName(gdb *gorm.DB) string {
	sqlDB, err := gdb.DB()
	if err != nil {
		return ""
	}
	var name string
	// Portable approach: SELECT current_database();
	_ = queryRow(sqlDB, `SELECT current_database()`).Scan(&name)
	return name
}

func queryRow(db *sql.DB, q string, args ...any) *sql.Row {
	return db.QueryRow(q, args...)
}
