package domain

import "gorm.io/datatypes"

type Item struct {
	ID          string            `json:"id"           db:"id"           gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID    string            `json:"tenant_id"    db:"tenant_id"    gorm:"type:uuid;index"`
	CategoryID  string            `json:"category_id"  db:"category_id"  gorm:"type:uuid;index"`
	Name        string            `json:"name"         db:"name"         gorm:"not null"`
	Description *string           `json:"description,omitempty" db:"description"`
	Price       int64             `json:"price"        db:"price"`
	PhotoURL    *string           `json:"photo_url,omitempty" db:"photo_url"`
	Flags       datatypes.JSONMap `json:"flags,omitempty"     db:"flags"     gorm:"type:jsonb"`
	IsActive    bool              `json:"is_active"    db:"is_active"     gorm:"default:true;index"`
}
