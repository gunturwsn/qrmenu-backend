package domain

type Table struct {
	ID       string `json:"id"         db:"id"         gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID string `json:"tenant_id"  db:"tenant_id"  gorm:"type:uuid;index"`
	Code     string `json:"code"       db:"code"       gorm:"index"`                 // boleh di-unique per-tenant bila perlu
	Name     string `json:"name"       db:"name"`
	Token    string `json:"token"      db:"token"      gorm:"uniqueIndex;not null"` // dipakai di /api/v1/table/{token}
	IsActive bool   `json:"is_active"  db:"is_active"  gorm:"default:true;index"`
}
