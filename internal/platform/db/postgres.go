package db

import (
	"database/sql"
	"fmt"
	"log"

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

	// Optional: ping & set timezone pada session (bagus untuk verifikasi)
	if sqlDB, err := gdb.DB(); err == nil {
		_ = sqlDB.Ping()
	}
	if err := gdb.Exec(`SET TIME ZONE 'Asia/Jakarta'`).Error; err != nil {
		log.Printf("warn: failed to SET TIME ZONE: %v", err)
	}

	return gdb
}

// Pastikan extension yang dibutuhkan ada (sekali jalan aman walau berulang)
func EnsureExtensions(db *gorm.DB) {
	_ = db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error
	// Jika di DB kamu pakai uuid-ossp: aktifkan juga
	// _ = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}

func AutoMigrate(db *gorm.DB) {
	EnsureExtensions(db)

	// Urutan aman (FK mengacu entity yang sudah ada)
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

	// (Opsional) set default timezone pada level database (berlaku untuk sesi baru)
	// Aman dipanggil berulang
	if derr := setDatabaseTimezone(db, "Asia/Jakarta"); derr != nil {
		log.Printf("warn: failed to set DB default timezone: %v", derr)
	}
}

func setDatabaseTimezone(db *gorm.DB, tz string) error {
	// Mengatur default parameter untuk database ini, bukan hanya session saat ini
	return db.Exec(fmt.Sprintf(`ALTER DATABASE %s SET TIME ZONE '%s'`, currentDBName(db), tz)).Error
}

// currentDBName membaca nama DB aktif dari koneksi *sql.DB.
func currentDBName(gdb *gorm.DB) string {
	sqlDB, err := gdb.DB()
	if err != nil {
		return ""
	}
	var name string
	// cara portable: SELECT current_database();
	_ = queryRow(sqlDB, `SELECT current_database()`).Scan(&name)
	return name
}

func queryRow(db *sql.DB, q string, args ...any) *sql.Row {
	return db.QueryRow(q, args...)
}
