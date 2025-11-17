package repository

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

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
	if err != nil {
		logging.RepoError("CategoryRepository.List", "query failed", "query_failed", err, "tenant_id", tenantID)
		return nil, err
	}
	logging.RepoInfo("CategoryRepository.List", "categories listed", "categories_listed", "tenant_id", tenantID, "count", len(xs))
	return xs, err
}

func (r *categoryRepo) Create(c *domain.Category) error {
	if err := r.db.Create(c).Error; err != nil {
		logging.RepoError("CategoryRepository.Create", "insert failed", "insert_failed", err, "tenant_id", c.TenantID)
		return err
	}
	logging.RepoInfo("CategoryRepository.Create", "category created", "category_created", "tenant_id", c.TenantID, "category_id", c.ID)
	return nil
}

func (r *categoryRepo) Replace(c *domain.Category) error {
	if err := r.db.Where("id = ? AND tenant_id = ?", c.ID, c.TenantID).Updates(c).Error; err != nil {
		logging.RepoError("CategoryRepository.Replace", "update failed", "update_failed", err, "tenant_id", c.TenantID, "category_id", c.ID)
		return err
	}
	logging.RepoInfo("CategoryRepository.Replace", "category updated", "category_updated", "tenant_id", c.TenantID, "category_id", c.ID)
	return nil
}

func (r *categoryRepo) Patch(tenantID, id string, fields map[string]any) (*domain.Category, error) {
	if err := r.db.Model(&domain.Category{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).Updates(fields).Error; err != nil {
		logging.RepoError("CategoryRepository.Patch", "update failed", "update_failed", err, "tenant_id", tenantID, "category_id", id)
		return nil, err
	}
	logging.RepoInfo("CategoryRepository.Patch", "category patched", "category_patched", "tenant_id", tenantID, "category_id", id)
	return r.FindByID(tenantID, id)
}

func (r *categoryRepo) Delete(tenantID, id string) error {
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Category{}).Error; err != nil {
		logging.RepoError("CategoryRepository.Delete", "delete failed", "delete_failed", err, "tenant_id", tenantID, "category_id", id)
		return err
	}
	logging.RepoInfo("CategoryRepository.Delete", "category deleted", "category_deleted", "tenant_id", tenantID, "category_id", id)
	return nil
}

func (r *categoryRepo) FindByID(tenantID, id string) (*domain.Category, error) {
	var c domain.Category
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error; err != nil {
		logging.RepoError("CategoryRepository.FindByID", "query failed", "query_failed", err, "tenant_id", tenantID, "category_id", id)
		return nil, err
	}
	logging.RepoInfo("CategoryRepository.FindByID", "category found", "category_found", "tenant_id", tenantID, "category_id", id)
	return &c, nil
}
