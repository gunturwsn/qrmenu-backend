package repository

import (
	"errors"

	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

	"gorm.io/gorm"
)

type TableRepository interface {
	ResolveByToken(token string) (*domain.Tenant, *domain.Table, error)
}

type tableRepo struct{ db *gorm.DB }

func NewTableRepository(db *gorm.DB) TableRepository { return &tableRepo{db} }

func (r *tableRepo) ResolveByToken(token string) (*domain.Tenant, *domain.Table, error) {
	var tb domain.Table
	if err := r.db.Where("token = ? AND is_active = TRUE", token).First(&tb).Error; err != nil {
		logging.RepoError("TableRepository.ResolveByToken", "table lookup failed", "table_lookup_failed", err, "token", token)
		return nil, nil, err
	}
	var tn domain.Tenant
	if err := r.db.First(&tn, "id = ?", tb.TenantID).Error; err != nil {
		logging.RepoError("TableRepository.ResolveByToken", "tenant lookup failed", "tenant_lookup_failed", err, "tenant_id", tb.TenantID, "token", token)
		return nil, nil, err
	}
	if !tb.IsActive {
		logging.RepoError("TableRepository.ResolveByToken", "table inactive", "table_inactive", errors.New("table inactive"), "token", token)
		return nil, nil, errors.New("table inactive")
	}
	logging.RepoInfo("TableRepository.ResolveByToken", "table resolved", "table_resolved", "token", token, "tenant_id", tn.ID)
	return &tn, &tb, nil
}
