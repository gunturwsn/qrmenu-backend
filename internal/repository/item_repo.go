package repository

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

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
	if err != nil {
		logging.RepoError("ItemRepository.List", "query failed", "query_failed", err, "tenant_id", tenantID, "category_id", categoryID)
		return nil, err
	}
	logging.RepoInfo("ItemRepository.List", "items listed", "items_listed", "tenant_id", tenantID, "category_id", categoryID, "count", len(xs))
	return xs, err
}

func (r *itemRepo) Create(i *domain.Item) error {
	if err := r.db.Create(i).Error; err != nil {
		logging.RepoError("ItemRepository.Create", "insert failed", "insert_failed", err, "tenant_id", i.TenantID)
		return err
	}
	logging.RepoInfo("ItemRepository.Create", "item created", "item_created", "tenant_id", i.TenantID, "item_id", i.ID)
	return nil
}

func (r *itemRepo) Replace(i *domain.Item) error {
	if err := r.db.Where("id = ? AND tenant_id = ?", i.ID, i.TenantID).Updates(i).Error; err != nil {
		logging.RepoError("ItemRepository.Replace", "update failed", "update_failed", err, "tenant_id", i.TenantID, "item_id", i.ID)
		return err
	}
	logging.RepoInfo("ItemRepository.Replace", "item updated", "item_updated", "tenant_id", i.TenantID, "item_id", i.ID)
	return nil
}

func (r *itemRepo) Patch(tenantID, id string, fields map[string]any) (*domain.Item, error) {
	if err := r.db.Model(&domain.Item{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).Updates(fields).Error; err != nil {
		logging.RepoError("ItemRepository.Patch", "update failed", "update_failed", err, "tenant_id", tenantID, "item_id", id)
		return nil, err
	}
	logging.RepoInfo("ItemRepository.Patch", "item patched", "item_patched", "tenant_id", tenantID, "item_id", id)
	return r.FindByID(tenantID, id)
}

func (r *itemRepo) Delete(tenantID, id string) error {
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Item{}).Error; err != nil {
		logging.RepoError("ItemRepository.Delete", "delete failed", "delete_failed", err, "tenant_id", tenantID, "item_id", id)
		return err
	}
	logging.RepoInfo("ItemRepository.Delete", "item deleted", "item_deleted", "tenant_id", tenantID, "item_id", id)
	return nil
}

func (r *itemRepo) FindByID(tenantID, id string) (*domain.Item, error) {
	var m domain.Item
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&m).Error; err != nil {
		logging.RepoError("ItemRepository.FindByID", "query failed", "query_failed", err, "tenant_id", tenantID, "item_id", id)
		return nil, err
	}
	logging.RepoInfo("ItemRepository.FindByID", "item found", "item_found", "tenant_id", tenantID, "item_id", id)
	return &m, nil
}

func (r *itemRepo) ToggleActive(tenantID, id string, isActive bool) (*domain.Item, error) {
	if err := r.db.Model(&domain.Item{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("is_active", isActive).Error; err != nil {
		logging.RepoError("ItemRepository.ToggleActive", "update failed", "update_failed", err, "tenant_id", tenantID, "item_id", id, "is_active", isActive)
		return nil, err
	}
	logging.RepoInfo("ItemRepository.ToggleActive", "item toggled", "item_toggled", "tenant_id", tenantID, "item_id", id, "is_active", isActive)
	return r.FindByID(tenantID, id)
}
