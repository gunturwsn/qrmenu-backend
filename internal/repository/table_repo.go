package repository

import (
	"errors"
	"qrmenu/internal/domain"

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
		return nil, nil, err
	}
	var tn domain.Tenant
	if err := r.db.First(&tn, "id = ?", tb.TenantID).Error; err != nil {
		return nil, nil, err
	}
	if !tb.IsActive {
		return nil, nil, errors.New("table inactive")
	}
	return &tn, &tb, nil
}
