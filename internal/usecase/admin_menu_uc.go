package usecase

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/repository"
)

type AdminMenuUC struct {
	catRepo  repository.CategoryRepository
	itemRepo repository.ItemRepository
	optRepo  repository.OptionRepository
}

func NewAdminMenuUC(cat repository.CategoryRepository, it repository.ItemRepository, op repository.OptionRepository) *AdminMenuUC {
	return &AdminMenuUC{catRepo: cat, itemRepo: it, optRepo: op}
}

// ===== Categories
func (u *AdminMenuUC) ListCategories(tenantID string) ([]domain.Category, error) {
	logging.UsecaseInfo("AdminMenu.ListCategories", "listing categories", "categories_list_requested", "tenant_id", tenantID)
	xs, err := u.catRepo.List(tenantID)
	if err != nil {
		logging.UsecaseError("AdminMenu.ListCategories", "repository error", "categories_list_failed", err, "tenant_id", tenantID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ListCategories", "categories loaded", "categories_listed", "tenant_id", tenantID, "count", len(xs))
	return xs, nil
}
func (u *AdminMenuUC) CreateCategory(tenantID string, body map[string]any) (*domain.Category, error) {
	logging.UsecaseInfo("AdminMenu.CreateCategory", "creating category", "category_create_requested", "tenant_id", tenantID)
	c := &domain.Category{TenantID: tenantID, Name: body["name"].(string)}
	if v, ok := body["sort"].(float64); ok {
		c.Sort = int(v)
	}
	if v, ok := body["is_active"].(bool); ok {
		c.IsActive = v
	}
	if err := u.catRepo.Create(c); err != nil {
		logging.UsecaseError("AdminMenu.CreateCategory", "repository error", "category_create_failed", err, "tenant_id", tenantID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.CreateCategory", "category created", "category_created", "tenant_id", tenantID, "category_id", c.ID)
	return c, nil
}
func (u *AdminMenuUC) ReplaceCategory(tenantID, id string, body map[string]any) (*domain.Category, error) {
	logging.UsecaseInfo("AdminMenu.ReplaceCategory", "replacing category", "category_replace_requested", "tenant_id", tenantID, "category_id", id)
	c := &domain.Category{ID: id, TenantID: tenantID, Name: body["name"].(string)}
	if v, ok := body["sort"].(float64); ok {
		c.Sort = int(v)
	}
	if v, ok := body["is_active"].(bool); ok {
		c.IsActive = v
	}
	if err := u.catRepo.Replace(c); err != nil {
		logging.UsecaseError("AdminMenu.ReplaceCategory", "repository error", "category_replace_failed", err, "tenant_id", tenantID, "category_id", id)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ReplaceCategory", "category updated", "category_replaced", "tenant_id", tenantID, "category_id", id)
	return c, nil
}
func (u *AdminMenuUC) PatchCategory(tenantID, id string, body map[string]any) (*domain.Category, error) {
	logging.UsecaseInfo("AdminMenu.PatchCategory", "patching category", "category_patch_requested", "tenant_id", tenantID, "category_id", id)
	obj, err := u.catRepo.Patch(tenantID, id, body)
	if err != nil {
		logging.UsecaseError("AdminMenu.PatchCategory", "repository error", "category_patch_failed", err, "tenant_id", tenantID, "category_id", id)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.PatchCategory", "category patched", "category_patched", "tenant_id", tenantID, "category_id", id)
	return obj, nil
}
func (u *AdminMenuUC) DeleteCategory(tenantID, id string) error {
	logging.UsecaseInfo("AdminMenu.DeleteCategory", "deleting category", "category_delete_requested", "tenant_id", tenantID, "category_id", id)
	if err := u.catRepo.Delete(tenantID, id); err != nil {
		logging.UsecaseError("AdminMenu.DeleteCategory", "repository error", "category_delete_failed", err, "tenant_id", tenantID, "category_id", id)
		return err
	}
	logging.UsecaseInfo("AdminMenu.DeleteCategory", "category deleted", "category_deleted", "tenant_id", tenantID, "category_id", id)
	return nil
}

// ===== Items
func (u *AdminMenuUC) ListItems(tenantID, categoryID string) ([]domain.Item, error) {
	logging.UsecaseInfo("AdminMenu.ListItems", "listing items", "items_list_requested", "tenant_id", tenantID, "category_id", categoryID)
	xs, err := u.itemRepo.List(tenantID, categoryID)
	if err != nil {
		logging.UsecaseError("AdminMenu.ListItems", "repository error", "items_list_failed", err, "tenant_id", tenantID, "category_id", categoryID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ListItems", "items loaded", "items_listed", "tenant_id", tenantID, "category_id", categoryID, "count", len(xs))
	return xs, nil
}
func (u *AdminMenuUC) CreateItem(tenantID string, body map[string]any) (*domain.Item, error) {
	logging.UsecaseInfo("AdminMenu.CreateItem", "creating item", "item_create_requested", "tenant_id", tenantID)
	i := &domain.Item{TenantID: tenantID}
	if v, ok := body["category_id"].(string); ok {
		i.CategoryID = v
	}
	if v, ok := body["name"].(string); ok {
		i.Name = v
	}
	if v, ok := body["description"].(string); ok {
		i.Description = &v
	}
	if v, ok := body["price"].(float64); ok {
		i.Price = int64(v)
	}
	if v, ok := body["photo_url"].(string); ok {
		i.PhotoURL = &v
	}
	if v, ok := body["is_active"].(bool); ok {
		i.IsActive = v
	}
	if err := u.itemRepo.Create(i); err != nil {
		logging.UsecaseError("AdminMenu.CreateItem", "repository error", "item_create_failed", err, "tenant_id", tenantID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.CreateItem", "item created", "item_created", "tenant_id", tenantID, "item_id", i.ID)
	return i, nil
}
func (u *AdminMenuUC) ReplaceItem(tenantID, id string, body map[string]any) (*domain.Item, error) {
	logging.UsecaseInfo("AdminMenu.ReplaceItem", "replacing item", "item_replace_requested", "tenant_id", tenantID, "item_id", id)
	i := &domain.Item{ID: id, TenantID: tenantID}
	if v, ok := body["category_id"].(string); ok {
		i.CategoryID = v
	}
	if v, ok := body["name"].(string); ok {
		i.Name = v
	}
	if v, ok := body["description"].(string); ok {
		i.Description = &v
	} else {
		i.Description = nil
	}
	if v, ok := body["price"].(float64); ok {
		i.Price = int64(v)
	}
	if v, ok := body["photo_url"].(string); ok {
		i.PhotoURL = &v
	} else {
		i.PhotoURL = nil
	}
	if v, ok := body["is_active"].(bool); ok {
		i.IsActive = v
	}
	if err := u.itemRepo.Replace(i); err != nil {
		logging.UsecaseError("AdminMenu.ReplaceItem", "repository error", "item_replace_failed", err, "tenant_id", tenantID, "item_id", id)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ReplaceItem", "item updated", "item_replaced", "tenant_id", tenantID, "item_id", id)
	return i, nil
}
func (u *AdminMenuUC) PatchItem(tenantID, id string, body map[string]any) (*domain.Item, error) {
	logging.UsecaseInfo("AdminMenu.PatchItem", "patching item", "item_patch_requested", "tenant_id", tenantID, "item_id", id)
	obj, err := u.itemRepo.Patch(tenantID, id, body)
	if err != nil {
		logging.UsecaseError("AdminMenu.PatchItem", "repository error", "item_patch_failed", err, "tenant_id", tenantID, "item_id", id)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.PatchItem", "item patched", "item_patched", "tenant_id", tenantID, "item_id", id)
	return obj, nil
}
func (u *AdminMenuUC) DeleteItem(tenantID, id string) error {
	logging.UsecaseInfo("AdminMenu.DeleteItem", "deleting item", "item_delete_requested", "tenant_id", tenantID, "item_id", id)
	if err := u.itemRepo.Delete(tenantID, id); err != nil {
		logging.UsecaseError("AdminMenu.DeleteItem", "repository error", "item_delete_failed", err, "tenant_id", tenantID, "item_id", id)
		return err
	}
	logging.UsecaseInfo("AdminMenu.DeleteItem", "item deleted", "item_deleted", "tenant_id", tenantID, "item_id", id)
	return nil
}
func (u *AdminMenuUC) ToggleOOS(tenantID, id string, isActive bool) (*domain.Item, error) {
	logging.UsecaseInfo("AdminMenu.ToggleOOS", "toggling item activity", "item_toggle_requested", "tenant_id", tenantID, "item_id", id, "is_active", isActive)
	obj, err := u.itemRepo.ToggleActive(tenantID, id, isActive)
	if err != nil {
		logging.UsecaseError("AdminMenu.ToggleOOS", "repository error", "item_toggle_failed", err, "tenant_id", tenantID, "item_id", id, "is_active", isActive)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ToggleOOS", "item toggled", "item_toggled", "tenant_id", tenantID, "item_id", id, "is_active", isActive)
	return obj, nil
}

// ===== Options
func (u *AdminMenuUC) ListItemOptions(itemID, tenantID string) ([]domain.ItemOption, error) {
	logging.UsecaseInfo("AdminMenu.ListItemOptions", "listing options", "options_list_requested", "tenant_id", tenantID, "item_id", itemID)
	xs, err := u.optRepo.ListItemOptions(itemID, tenantID)
	if err != nil {
		logging.UsecaseError("AdminMenu.ListItemOptions", "repository error", "options_list_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ListItemOptions", "options loaded", "options_listed", "tenant_id", tenantID, "item_id", itemID, "count", len(xs))
	return xs, nil
}
func (u *AdminMenuUC) CreateItemOption(itemID, tenantID string, body map[string]any) (*domain.ItemOption, error) {
	logging.UsecaseInfo("AdminMenu.CreateItemOption", "creating option", "option_create_requested", "tenant_id", tenantID, "item_id", itemID)
	o := &domain.ItemOption{Name: body["name"].(string), Type: body["type"].(string)}
	if v, ok := body["required"].(bool); ok {
		o.Required = v
	}
	if err := u.optRepo.CreateItemOption(itemID, tenantID, o); err != nil {
		logging.UsecaseError("AdminMenu.CreateItemOption", "repository error", "option_create_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.CreateItemOption", "option created", "option_created", "tenant_id", tenantID, "item_id", itemID, "option_id", o.ID)
	return o, nil
}
func (u *AdminMenuUC) ListOptionValues(optionID, tenantID string) ([]domain.ItemOptionValue, error) {
	logging.UsecaseInfo("AdminMenu.ListOptionValues", "listing option values", "option_values_list_requested", "tenant_id", tenantID, "option_id", optionID)
	xs, err := u.optRepo.ListOptionValues(optionID, tenantID)
	if err != nil {
		logging.UsecaseError("AdminMenu.ListOptionValues", "repository error", "option_values_list_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.ListOptionValues", "option values loaded", "option_values_listed", "tenant_id", tenantID, "option_id", optionID, "count", len(xs))
	return xs, nil
}
func (u *AdminMenuUC) CreateOptionValue(optionID, tenantID string, body map[string]any) (*domain.ItemOptionValue, error) {
	logging.UsecaseInfo("AdminMenu.CreateOptionValue", "creating option value", "option_value_create_requested", "tenant_id", tenantID, "option_id", optionID)
	v := &domain.ItemOptionValue{Label: body["label"].(string)}
	if dp, ok := body["delta_price"].(float64); ok {
		v.DeltaPrice = int64(dp)
	}
	if err := u.optRepo.CreateOptionValue(optionID, tenantID, v); err != nil {
		logging.UsecaseError("AdminMenu.CreateOptionValue", "repository error", "option_value_create_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return nil, err
	}
	logging.UsecaseInfo("AdminMenu.CreateOptionValue", "option value created", "option_value_created", "tenant_id", tenantID, "option_id", optionID, "value_id", v.ID)
	return v, nil
}

// ===== Tables (QR)
func (u *AdminMenuUC) GenerateTableQR(tableID, tenantID string) (string, error) {
	logging.UsecaseInfo("AdminMenu.GenerateTableQR", "generating qr", "qr_generate_requested", "tenant_id", tenantID, "table_id", tableID)
	// TODO: implement real QR generation / storage
	url := "https://example.com/qr/" + tableID
	logging.UsecaseInfo("AdminMenu.GenerateTableQR", "qr generated", "qr_generated", "tenant_id", tenantID, "table_id", tableID)
	return url, nil
}
