package repository

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

	"gorm.io/gorm"
)

type MenuQuery interface {
	GetMenuByTenantCode(code string) (*domain.MenuResponse, error)
}

type menuQuery struct{ db *gorm.DB }

func NewMenuQuery(db *gorm.DB) MenuQuery { return &menuQuery{db} }

// GetMenuByTenantCode returns a menu response based on the given tenant code.
// It first finds the tenant based on the given code, then retrieves the categories and items
// for the tenant. If the tenant is not found, it returns an error.
// If there is an error during the database query, it also returns an error.
func (q *menuQuery) GetMenuByTenantCode(code string) (*domain.MenuResponse, error) {
	var t domain.Tenant
	err := q.db.Where("code = ?", code).First(&t).Error
	if err != nil {
		logging.RepoError("MenuQuery.GetMenuByTenantCode", "tenant lookup failed", "tenant_lookup_failed", err, "tenant_code", code)
		return nil, err
	}
	var cats []domain.Category
	if err := q.db.Where("tenant_id = ? AND is_active = TRUE", t.ID).
		Order("sort ASC, name ASC").Find(&cats).Error; err != nil {
		logging.RepoError("MenuQuery.GetMenuByTenantCode", "categories lookup failed", "categories_query_failed", err, "tenant_id", t.ID)
		return nil, err
	}
	var items []domain.Item
	if err := q.db.Where("tenant_id = ? AND is_active = TRUE", t.ID).
		Order("name ASC").Find(&items).Error; err != nil {
		logging.RepoError("MenuQuery.GetMenuByTenantCode", "items lookup failed", "items_query_failed", err, "tenant_id", t.ID)
		return nil, err
	}
	logging.RepoInfo("MenuQuery.GetMenuByTenantCode", "menu loaded", "menu_loaded", "tenant_code", code, "categories", len(cats), "items", len(items))
	return &domain.MenuResponse{
		Tenant:     t.Code,
		Categories: cats,
		Items:      items,
	}, nil
}
