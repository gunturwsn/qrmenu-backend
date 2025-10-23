package domain

import "time"

type AdminUser struct {
	ID           string    `json:"id"          db:"id"            gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     string    `json:"tenant_id"   db:"tenant_id"     gorm:"type:uuid;index"`
	Email        string    `json:"email"       db:"email"         gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-"           db:"password_hash" gorm:"not null"`
	Name         string    `json:"name"        db:"name"`
	Role         string    `json:"role"        db:"role"          gorm:"default:'staff'"` // owner|staff
	IsActive     bool      `json:"is_active"   db:"is_active"     gorm:"default:true;index"`
	CreatedAt    time.Time `json:"created_at"  db:"created_at"    gorm:"autoCreateTime"`
}
