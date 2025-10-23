package domain

import (
	"time"

	"gorm.io/datatypes"
)

type OrderStatus string
const (
	OrderWaiting    OrderStatus = "waiting"
	OrderProcessing OrderStatus = "processing"
	OrderDelivering OrderStatus = "delivering"
	OrderDone       OrderStatus = "done"
	OrderCanceled   OrderStatus = "canceled"
)

type PaidStatus string
const (
	Unpaid PaidStatus = "unpaid"
	Paid   PaidStatus = "paid"
)

type Order struct {
	ID           string      `json:"id"               db:"id"               gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     string      `json:"tenant_id"        db:"tenant_id"        gorm:"type:uuid;index"`
	TableID      string      `json:"table_id"         db:"table_id"         gorm:"type:uuid;index"`
	GuestSession string      `json:"guest_session_id,omitempty" db:"guest_session_id"`
	Note         *string     `json:"note,omitempty"   db:"note"`
	Status       OrderStatus `json:"status"           db:"status"           gorm:"type:text;default:'waiting';index"`
	PaidStatus   PaidStatus  `json:"paid_status"      db:"paid_status"      gorm:"type:text;default:'unpaid';index"`
	CreatedAt    time.Time   `json:"created_at"       db:"created_at"       gorm:"autoCreateTime"`

	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID        string            `json:"id"        db:"id"        gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   string            `json:"order_id"  db:"order_id"  gorm:"type:uuid;index"`
	ItemID    string            `json:"item_id"   db:"item_id"   gorm:"type:uuid;index"`
	Name      string            `json:"name"      db:"name"`
	Qty       int               `json:"qty"       db:"qty"`
	UnitPrice int64             `json:"unit_price" db:"unit_price"`
	Options   datatypes.JSONMap `json:"options,omitempty" db:"options" gorm:"type:jsonb"`
}
