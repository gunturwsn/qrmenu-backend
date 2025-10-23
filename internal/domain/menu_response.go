package domain

type MenuResponse struct {
	Tenant     string     `json:"tenant"`
	Categories []Category `json:"categories"`
	Items      []Item     `json:"items"`
}
