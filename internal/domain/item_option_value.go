package domain

type ItemOptionValue struct {
	ID         string `json:"id"          db:"id"          gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OptionID   string `json:"option_id"   db:"option_id"   gorm:"type:uuid;index"`
	Label      string `json:"label"       db:"label"       gorm:"not null"`
	DeltaPrice int64  `json:"delta_price" db:"delta_price" gorm:"default:0"`
}
