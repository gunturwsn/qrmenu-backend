package repository

import (
	"qrmenu/internal/domain"

	"gorm.io/gorm"
)

type ItemRepository interface {
	List(tenantID string, categoryID string) ([]domain.Item, error)
	Create(i *domain.Item) error
	Replace(i *domain.Item) error
	Patch(tenantID, id string, fields map[string]any) (*domain.Item, error)
	Delete(tenantID, id string) error
	FindByID(tenantID, id string) (*domain.Item, error)
	ToggleActive(tenantID, id string, isActive bool) (*domain.Item, error)
}

type itemRepo struct{ db *gorm.DB }

func NewItemRepository(db *gorm.DB) ItemRepository { return &itemRepo{db} }

func (r *itemRepo) List(tenantID string, categoryID string) ([]domain.Item, error) {
	q := r.db.Where("tenant_id = ?", tenantID)
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	var xs []domain.Item
	err := q.Order("name ASC").Find(&xs).Error
	return xs, err
}

func (r *itemRepo) Create(i *domain.Item) error { return r.db.Create(i).Error }

func (r *itemRepo) Replace(i *domain.Item) error {
	return r.db.Where("id = ? AND tenant_id = ?", i.ID, i.TenantID).Updates(i).Error
}

func (r *itemRepo) Patch(tenantID, id string, fields map[string]any) (*domain.Item, error) {
	if err := r.db.Model(&domain.Item{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).Updates(fields).Error; err != nil {
		return nil, err
	}
	return r.FindByID(tenantID, id)
}

func (r *itemRepo) Delete(tenantID, id string) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Item{}).Error
}

func (r *itemRepo) FindByID(tenantID, id string) (*domain.Item, error) {
	var m domain.Item
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *itemRepo) ToggleActive(tenantID, id string, isActive bool) (*domain.Item, error) {
	if err := r.db.Model(&domain.Item{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("is_active", isActive).Error; err != nil {
		return nil, err
	}
	return r.FindByID(tenantID, id)
}
