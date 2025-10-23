package domain

type OrderItemCreate struct {
	ItemID  string         `json:"item_id"`
	Qty     int            `json:"qty"`
	Options map[string]any `json:"options,omitempty"`
}

type OrderCreateRequest struct {
	Tenant        string             `json:"tenant"`
	TableToken    string             `json:"table_token"`
	GuestSession  string             `json:"guest_session_id"`
	Note          *string            `json:"note,omitempty"`
	Items         []OrderItemCreate  `json:"items"`
}
