package domain

import (
	"time"

	"gorm.io/datatypes"
)

type Tenant struct {
	ID        string             `json:"id"       db:"id"       gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      string             `json:"code"     db:"code"     gorm:"uniqueIndex;not null"`
	Name      string             `json:"name"     db:"name"`
	LogoURL   *string            `json:"logo_url,omitempty" db:"logo_url"`
	Theme     datatypes.JSONMap  `json:"theme,omitempty"    db:"theme"    gorm:"type:jsonb"`
	CreatedAt time.Time          `json:"created_at"         db:"created_at" gorm:"autoCreateTime"`
}
