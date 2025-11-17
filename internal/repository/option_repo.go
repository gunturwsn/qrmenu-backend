package repository

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

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
	// Guard tenant ownership via join to items table
	var xs []domain.ItemOption
	err := r.db.Table("item_options io").
		Select("io.*").
		Joins("JOIN items i ON i.id = io.item_id AND i.tenant_id = ?", tenantID).
		Where("io.item_id = ?", itemID).
		Order("io.name ASC").Scan(&xs).Error
	if err != nil {
		logging.RepoError("OptionRepository.ListItemOptions", "query failed", "query_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return nil, err
	}
	logging.RepoInfo("OptionRepository.ListItemOptions", "options listed", "options_listed", "tenant_id", tenantID, "item_id", itemID, "count", len(xs))
	return xs, err
}

func (r *optionRepo) CreateItemOption(itemID, tenantID string, opt *domain.ItemOption) error {
	// Verify the item belongs to the tenant
	var cnt int64
	if err := r.db.Model(&domain.Item{}).
		Where("id = ? AND tenant_id = ?", itemID, tenantID).Count(&cnt).Error; err != nil {
		logging.RepoError("OptionRepository.CreateItemOption", "item validation failed", "item_validation_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return err
	}
	if cnt == 0 {
		return gorm.ErrRecordNotFound
	}
	opt.ItemID = itemID
	if err := r.db.Create(opt).Error; err != nil {
		logging.RepoError("OptionRepository.CreateItemOption", "insert failed", "insert_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return err
	}
	logging.RepoInfo("OptionRepository.CreateItemOption", "option created", "option_created", "tenant_id", tenantID, "item_id", itemID, "option_id", opt.ID)
	return nil
}

func (r *optionRepo) ListOptionValues(optionID, tenantID string) ([]domain.ItemOptionValue, error) {
	var xs []domain.ItemOptionValue
	err := r.db.Table("item_option_values v").
		Select("v.*").
		Joins("JOIN item_options o ON o.id = v.option_id").
		Joins("JOIN items i ON i.id = o.item_id AND i.tenant_id = ?", tenantID).
		Where("v.option_id = ?", optionID).
		Order("v.label ASC").Scan(&xs).Error
	if err != nil {
		logging.RepoError("OptionRepository.ListOptionValues", "query failed", "query_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return nil, err
	}
	logging.RepoInfo("OptionRepository.ListOptionValues", "option values listed", "option_values_listed", "tenant_id", tenantID, "option_id", optionID, "count", len(xs))
	return xs, err
}

func (r *optionRepo) CreateOptionValue(optionID, tenantID string, v *domain.ItemOptionValue) error {
	// Ensure option → item → tenant ownership aligns
	var cnt int64
	if err := r.db.Table("item_options o").
		Joins("JOIN items i ON i.id = o.item_id").
		Where("o.id = ? AND i.tenant_id = ?", optionID, tenantID).
		Count(&cnt).Error; err != nil {
		logging.RepoError("OptionRepository.CreateOptionValue", "option validation failed", "option_validation_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return err
	}
	if cnt == 0 {
		return gorm.ErrRecordNotFound
	}
	v.OptionID = optionID
	if err := r.db.Create(v).Error; err != nil {
		logging.RepoError("OptionRepository.CreateOptionValue", "insert failed", "insert_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return err
	}
	logging.RepoInfo("OptionRepository.CreateOptionValue", "option value created", "option_value_created", "tenant_id", tenantID, "option_id", optionID, "value_id", v.ID)
	return nil
}
