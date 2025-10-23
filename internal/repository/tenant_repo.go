package repository

import (
	"qrmenu/internal/domain"

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
		return &domain.Tenant{}, err
	}
	return &t, nil
}

func (r *tenantRepo) Create(t *domain.Tenant) error {
	return r.db.Create(t).Error
}

