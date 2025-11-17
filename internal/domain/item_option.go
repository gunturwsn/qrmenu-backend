package domain

type ItemOption struct {
	ID       string `json:"id"       db:"id"       gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ItemID   string `json:"item_id"  db:"item_id"  gorm:"type:uuid;index"`
	Name     string `json:"name"     db:"name"     gorm:"not null"`
	Type     string `json:"type"     db:"type"     gorm:"not null"`
	Required bool   `json:"required" db:"required" gorm:"default:false"`
}
