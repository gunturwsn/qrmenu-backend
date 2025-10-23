package handler

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/usecase"
)

type adminMenuUCAdapter struct{ uc *usecase.AdminMenuUC }

func NewAdminMenuServiceFromUC(uc *usecase.AdminMenuUC) *adminMenuUCAdapter {
	return &adminMenuUCAdapter{uc: uc}
}

// Categories
func (a *adminMenuUCAdapter) ListCategories(tenantID string) ([]map[string]any, error) {
	cs, err := a.uc.ListCategories(tenantID)
	if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(cs))
	for _, c := range cs { out = append(out, mapCat(c)) }
	return out, nil
}
func (a *adminMenuUCAdapter) CreateCategory(tenantID string, body map[string]any) (map[string]any, error) {
	c, err := a.uc.CreateCategory(tenantID, body); if err != nil { return nil, err }
	return mapCat(*c), nil
}
func (a *adminMenuUCAdapter) ReplaceCategory(tenantID, id string, body map[string]any) (map[string]any, error) {
	c, err := a.uc.ReplaceCategory(tenantID, id, body); if err != nil { return nil, err }
	return mapCat(*c), nil
}
func (a *adminMenuUCAdapter) PatchCategory(tenantID, id string, body map[string]any) (map[string]any, error) {
	c, err := a.uc.PatchCategory(tenantID, id, body); if err != nil { return nil, err }
	return mapCat(*c), nil
}
func (a *adminMenuUCAdapter) DeleteCategory(tenantID, id string) error {
	return a.uc.DeleteCategory(tenantID, id)
}

// Items
func (a *adminMenuUCAdapter) ListItems(tenantID, categoryID string) ([]map[string]any, error) {
	items, err := a.uc.ListItems(tenantID, categoryID)
	if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(items))
	for _, m := range items { out = append(out, mapItem(m)) }
	return out, nil
}
func (a *adminMenuUCAdapter) CreateItem(tenantID string, body map[string]any) (map[string]any, error) {
	i, err := a.uc.CreateItem(tenantID, body); if err != nil { return nil, err }
	return mapItem(*i), nil
}
func (a *adminMenuUCAdapter) ReplaceItem(tenantID, id string, body map[string]any) (map[string]any, error) {
	i, err := a.uc.ReplaceItem(tenantID, id, body); if err != nil { return nil, err }
	return mapItem(*i), nil
}
func (a *adminMenuUCAdapter) PatchItem(tenantID, id string, body map[string]any) (map[string]any, error) {
	i, err := a.uc.PatchItem(tenantID, id, body); if err != nil { return nil, err }
	return mapItem(*i), nil
}
func (a *adminMenuUCAdapter) DeleteItem(tenantID, id string) error {
	return a.uc.DeleteItem(tenantID, id)
}
func (a *adminMenuUCAdapter) ToggleOOS(tenantID, id string, isActive bool) (map[string]any, error) {
	i, err := a.uc.ToggleOOS(tenantID, id, isActive); if err != nil { return nil, err }
	return mapItem(*i), nil
}

// Options
func (a *adminMenuUCAdapter) ListItemOptions(itemID, tenantID string) ([]map[string]any, error) {
	opts, err := a.uc.ListItemOptions(itemID, tenantID)
	if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(opts))
	for _, o := range opts { out = append(out, mapOption(o)) }
	return out, nil
}
func (a *adminMenuUCAdapter) CreateItemOption(itemID, tenantID string, body map[string]any) (map[string]any, error) {
	o, err := a.uc.CreateItemOption(itemID, tenantID, body); if err != nil { return nil, err }
	return mapOption(*o), nil
}
func (a *adminMenuUCAdapter) ListOptionValues(optionID, tenantID string) ([]map[string]any, error) {
	vals, err := a.uc.ListOptionValues(optionID, tenantID); if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(vals))
	for _, v := range vals { out = append(out, mapOptionValue(v)) }
	return out, nil
}
func (a *adminMenuUCAdapter) CreateOptionValue(optionID, tenantID string, body map[string]any) (map[string]any, error) {
	v, err := a.uc.CreateOptionValue(optionID, tenantID, body); if err != nil { return nil, err }
	return mapOptionValue(*v), nil
}

// Tables
func (a *adminMenuUCAdapter) GenerateTableQR(tableID, tenantID string) (string, error) {
	return a.uc.GenerateTableQR(tableID, tenantID)
}

// ----- mappers -----
func mapCat(c domain.Category) map[string]any {
	return map[string]any{
		"id": c.ID, "tenant_id": c.TenantID, "name": c.Name,
		"sort": c.Sort, "is_active": c.IsActive,
	}
}
func mapItem(i domain.Item) map[string]any {
	return map[string]any{
		"id": i.ID, "tenant_id": i.TenantID, "category_id": i.CategoryID,
		"name": i.Name, "description": i.Description, "price": i.Price,
		"photo_url": i.PhotoURL, "flags": i.Flags, "is_active": i.IsActive,
	}
}
func mapOption(o domain.ItemOption) map[string]any {
	return map[string]any{
		"id": o.ID, "item_id": o.ItemID, "name": o.Name, "type": o.Type, "required": o.Required,
	}
}
func mapOptionValue(v domain.ItemOptionValue) map[string]any {
	return map[string]any{
		"id": v.ID, "option_id": v.OptionID, "label": v.Label, "delta_price": v.DeltaPrice,
	}
}
