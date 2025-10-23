package repository

import (
	"qrmenu/internal/domain"

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
  return n, err
}

func (r *adminRepo) FindActiveByEmail(email string) (*domain.AdminUser, error) {
	var a domain.AdminUser
	if err := r.db.Where("email = ? AND is_active = true", email).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *adminRepo) CreateForTenant(a *domain.AdminUser) error {
	return r.db.Create(a).Error
}

func (r *adminRepo) Count() (int64, error) {
	var n int64
	err := r.db.Model(&domain.AdminUser{}).Where("is_active = ?", true).Count(&n).Error
	return n, err
}
