package domain

type Category struct {
	ID       string `json:"id"         db:"id"         gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID string `json:"tenant_id"  db:"tenant_id"  gorm:"type:uuid;index"`
	Name     string `json:"name"       db:"name"       gorm:"not null"`
	Sort     int    `json:"sort"       db:"sort"       gorm:"default:0"`
	IsActive bool   `json:"is_active"  db:"is_active"  gorm:"default:true;index"`
}
