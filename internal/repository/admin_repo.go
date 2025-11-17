package repository

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

	"gorm.io/gorm"
)

type AdminRepository interface {
	FindActiveByEmail(email string) (*domain.AdminUser, error)
	CreateForTenant(a *domain.AdminUser) error
	Count() (int64, error)
	CountActiveByTenant(tenantID string) (int64, error)
}

type adminRepo struct{ db *gorm.DB }

func NewAdminRepository(db *gorm.DB) AdminRepository { return &adminRepo{db: db} }

// repository
func (r *adminRepo) CountActiveByTenant(tenantID string) (int64, error) {
	var n int64
	err := r.db.Model(&domain.AdminUser{}).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Count(&n).Error
	if err != nil {
		logging.RepoError("AdminRepository.CountActiveByTenant", "count failed", "count_failed", err, "tenant_id", tenantID)
	}
	return n, err
}

func (r *adminRepo) FindActiveByEmail(email string) (*domain.AdminUser, error) {
	var a domain.AdminUser
	if err := r.db.Where("email = ? AND is_active = true", email).First(&a).Error; err != nil {
		logging.RepoError("AdminRepository.FindActiveByEmail", "query failed", "query_failed", err, "email", email)
		return nil, err
	}
	return &a, nil
}

func (r *adminRepo) CreateForTenant(a *domain.AdminUser) error {
	if err := r.db.Create(a).Error; err != nil {
		logging.RepoError("AdminRepository.CreateForTenant", "insert failed", "insert_failed", err, "tenant_id", a.TenantID, "email", a.Email)
		return err
	}
	logging.RepoInfo("AdminRepository.CreateForTenant", "admin created", "admin_created", "tenant_id", a.TenantID, "email", a.Email)
	return nil
}

func (r *adminRepo) Count() (int64, error) {
	var n int64
	err := r.db.Model(&domain.AdminUser{}).Where("is_active = ?", true).Count(&n).Error
	if err != nil {
		logging.RepoError("AdminRepository.Count", "count failed", "count_failed", err)
	}
	return n, err
}
