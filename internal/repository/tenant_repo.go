package repository

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

	"gorm.io/gorm"
)

type TenantRepository interface {
	FindByCode(code string) (*domain.Tenant, error)
	Create(t *domain.Tenant) error
}

type tenantRepo struct{ db *gorm.DB }

func NewTenantRepository(db *gorm.DB) TenantRepository { return &tenantRepo{db: db} }

func (r *tenantRepo) FindByCode(code string) (*domain.Tenant, error) {
	var t domain.Tenant
	if err := r.db.Where("code = ?", code).First(&t).Error; err != nil {
		logging.RepoError("TenantRepository.FindByCode", "query failed", "query_failed", err, "tenant_code", code)
		return &domain.Tenant{}, err
	}
	logging.RepoInfo("TenantRepository.FindByCode", "tenant found", "tenant_found", "tenant_code", code, "tenant_id", t.ID)
	return &t, nil
}

func (r *tenantRepo) Create(t *domain.Tenant) error {
	if err := r.db.Create(t).Error; err != nil {
		logging.RepoError("TenantRepository.Create", "insert failed", "insert_failed", err, "tenant_code", t.Code)
		return err
	}
	logging.RepoInfo("TenantRepository.Create", "tenant created", "tenant_created", "tenant_code", t.Code, "tenant_id", t.ID)
	return nil
}
