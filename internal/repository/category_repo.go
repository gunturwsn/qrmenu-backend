package repository

import (
	"qrmenu/internal/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	List(tenantID string) ([]domain.Category, error)
	Create(c *domain.Category) error
	Replace(c *domain.Category) error
	Patch(tenantID, id string, fields map[string]any) (*domain.Category, error)
	Delete(tenantID, id string) error
	FindByID(tenantID, id string) (*domain.Category, error)
}

type categoryRepo struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) CategoryRepository { return &categoryRepo{db} }

func (r *categoryRepo) List(tenantID string) ([]domain.Category, error) {
	var xs []domain.Category
	err := r.db.Where("tenant_id = ?", tenantID).Order("sort ASC, name ASC").Find(&xs).Error
	return xs, err
}

func (r *categoryRepo) Create(c *domain.Category) error { return r.db.Create(c).Error }

func (r *categoryRepo) Replace(c *domain.Category) error {
	return r.db.Where("id = ? AND tenant_id = ?", c.ID, c.TenantID).Updates(c).Error
}

func (r *categoryRepo) Patch(tenantID, id string, fields map[string]any) (*domain.Category, error) {
	if err := r.db.Model(&domain.Category{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).Updates(fields).Error; err != nil {
		return nil, err
	}
	return r.FindByID(tenantID, id)
}

func (r *categoryRepo) Delete(tenantID, id string) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Category{}).Error
}

func (r *categoryRepo) FindByID(tenantID, id string) (*domain.Category, error) {
	var c domain.Category
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
