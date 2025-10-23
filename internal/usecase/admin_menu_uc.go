package usecase

import (
	"qrmenu/internal/domain"
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
	return u.catRepo.List(tenantID)
}
func (u *AdminMenuUC) CreateCategory(tenantID string, body map[string]any) (*domain.Category, error) {
	c := &domain.Category{TenantID: tenantID, Name: body["name"].(string)}
	if v, ok := body["sort"].(float64); ok {
		c.Sort = int(v)
	}
	if v, ok := body["is_active"].(bool); ok {
		c.IsActive = v
	}
	return c, u.catRepo.Create(c)
}
func (u *AdminMenuUC) ReplaceCategory(tenantID, id string, body map[string]any) (*domain.Category, error) {
	c := &domain.Category{ID: id, TenantID: tenantID, Name: body["name"].(string)}
	if v, ok := body["sort"].(float64); ok {
		c.Sort = int(v)
	}
	if v, ok := body["is_active"].(bool); ok {
		c.IsActive = v
	}
	return c, u.catRepo.Replace(c)
}
func (u *AdminMenuUC) PatchCategory(tenantID, id string, body map[string]any) (*domain.Category, error) {
	return u.catRepo.Patch(tenantID, id, body)
}
func (u *AdminMenuUC) DeleteCategory(tenantID, id string) error {
	return u.catRepo.Delete(tenantID, id)
}

// ===== Items
func (u *AdminMenuUC) ListItems(tenantID, categoryID string) ([]domain.Item, error) {
	return u.itemRepo.List(tenantID, categoryID)
}
func (u *AdminMenuUC) CreateItem(tenantID string, body map[string]any) (*domain.Item, error) {
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
	return i, u.itemRepo.Create(i)
}
func (u *AdminMenuUC) ReplaceItem(tenantID, id string, body map[string]any) (*domain.Item, error) {
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
	return i, u.itemRepo.Replace(i)
}
func (u *AdminMenuUC) PatchItem(tenantID, id string, body map[string]any) (*domain.Item, error) {
	return u.itemRepo.Patch(tenantID, id, body)
}
func (u *AdminMenuUC) DeleteItem(tenantID, id string) error {
	return u.itemRepo.Delete(tenantID, id)
}
func (u *AdminMenuUC) ToggleOOS(tenantID, id string, isActive bool) (*domain.Item, error) {
	return u.itemRepo.ToggleActive(tenantID, id, isActive)
}

// ===== Options
func (u *AdminMenuUC) ListItemOptions(itemID, tenantID string) ([]domain.ItemOption, error) {
	return u.optRepo.ListItemOptions(itemID, tenantID)
}
func (u *AdminMenuUC) CreateItemOption(itemID, tenantID string, body map[string]any) (*domain.ItemOption, error) {
	o := &domain.ItemOption{Name: body["name"].(string), Type: body["type"].(string)}
	if v, ok := body["required"].(bool); ok {
		o.Required = v
	}
	return o, u.optRepo.CreateItemOption(itemID, tenantID, o)
}
func (u *AdminMenuUC) ListOptionValues(optionID, tenantID string) ([]domain.ItemOptionValue, error) {
	return u.optRepo.ListOptionValues(optionID, tenantID)
}
func (u *AdminMenuUC) CreateOptionValue(optionID, tenantID string, body map[string]any) (*domain.ItemOptionValue, error) {
	v := &domain.ItemOptionValue{Label: body["label"].(string)}
	if dp, ok := body["delta_price"].(float64); ok {
		v.DeltaPrice = int64(dp)
	}
	return v, u.optRepo.CreateOptionValue(optionID, tenantID, v)
}

// ===== Tables (QR)
func (u *AdminMenuUC) GenerateTableQR(tableID, tenantID string) (string, error) {
	// TODO: implement real QR generation / storage
	return "https://example.com/qr/" + tableID, nil
}
