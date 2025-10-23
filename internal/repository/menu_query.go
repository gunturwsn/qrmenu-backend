package repository

import (
	"qrmenu/internal/domain"

	"gorm.io/gorm"
)

type MenuQuery interface {
	GetMenuByTenantCode(code string) (*domain.MenuResponse, error)
}

type menuQuery struct{ db *gorm.DB }

func NewMenuQuery(db *gorm.DB) MenuQuery { return &menuQuery{db} }

func (q *menuQuery) GetMenuByTenantCode(code string) (*domain.MenuResponse, error) {
	var t domain.Tenant
	if err := q.db.Where("code = ?", code).First(&t).Error; err != nil {
		return nil, err
	}
	var cats []domain.Category
	if err := q.db.Where("tenant_id = ? AND is_active = TRUE", t.ID).
		Order("sort ASC, name ASC").Find(&cats).Error; err != nil {
		return nil, err
	}
	var items []domain.Item
	if err := q.db.Where("tenant_id = ? AND is_active = TRUE", t.ID).
		Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return &domain.MenuResponse{
		Tenant:     t.Code,
		Categories: cats,
		Items:      items,
	}, nil
}
