package repository

import (
	"qrmenu/internal/domain"

	"gorm.io/gorm"
)

type OptionRepository interface {
	ListItemOptions(itemID, tenantID string) ([]domain.ItemOption, error)
	CreateItemOption(itemID, tenantID string, opt *domain.ItemOption) error

	ListOptionValues(optionID, tenantID string) ([]domain.ItemOptionValue, error)
	CreateOptionValue(optionID, tenantID string, v *domain.ItemOptionValue) error
}

type optionRepo struct{ db *gorm.DB }

func NewOptionRepository(db *gorm.DB) OptionRepository { return &optionRepo{db} }

func (r *optionRepo) ListItemOptions(itemID, tenantID string) ([]domain.ItemOption, error) {
	// tenant guard via join ke item
	var xs []domain.ItemOption
	err := r.db.Table("item_options io").
		Select("io.*").
		Joins("JOIN items i ON i.id = io.item_id AND i.tenant_id = ?", tenantID).
		Where("io.item_id = ?", itemID).
		Order("io.name ASC").Scan(&xs).Error
	return xs, err
}

func (r *optionRepo) CreateItemOption(itemID, tenantID string, opt *domain.ItemOption) error {
	// validasi item milik tenant
	var cnt int64
	if err := r.db.Model(&domain.Item{}).
		Where("id = ? AND tenant_id = ?", itemID, tenantID).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 { return gorm.ErrRecordNotFound }
	opt.ItemID = itemID
	return r.db.Create(opt).Error
}

func (r *optionRepo) ListOptionValues(optionID, tenantID string) ([]domain.ItemOptionValue, error) {
	var xs []domain.ItemOptionValue
	err := r.db.Table("item_option_values v").
		Select("v.*").
		Joins("JOIN item_options o ON o.id = v.option_id").
		Joins("JOIN items i ON i.id = o.item_id AND i.tenant_id = ?", tenantID).
		Where("v.option_id = ?", optionID).
		Order("v.label ASC").Scan(&xs).Error
	return xs, err
}

func (r *optionRepo) CreateOptionValue(optionID, tenantID string, v *domain.ItemOptionValue) error {
	// pastikan option -> item -> tenant sesuai
	var cnt int64
	if err := r.db.Table("item_options o").
		Joins("JOIN items i ON i.id = o.item_id").
		Where("o.id = ? AND i.tenant_id = ?", optionID, tenantID).
		Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 { return gorm.ErrRecordNotFound }
	v.OptionID = optionID
	return r.db.Create(v).Error
}
