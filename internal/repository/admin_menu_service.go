package repository

import "qrmenu/internal/domain"

type AdminMenuService interface {
	// Categories
	ListCategories(tenantID string) ([]map[string]any, error)
	CreateCategory(tenantID string, body map[string]any) (map[string]any, error)
	ReplaceCategory(tenantID, id string, body map[string]any) (map[string]any, error)
	PatchCategory(tenantID, id string, body map[string]any) (map[string]any, error)
	DeleteCategory(tenantID, id string) error

	// Items
	ListItems(tenantID, categoryID string) ([]map[string]any, error)
	CreateItem(tenantID string, body map[string]any) (map[string]any, error)
	ReplaceItem(tenantID, id string, body map[string]any) (map[string]any, error)
	PatchItem(tenantID, id string, body map[string]any) (map[string]any, error)
	DeleteItem(tenantID, id string) error
	ToggleOOS(tenantID, id string, isActive bool) (map[string]any, error)

	// Options
	ListItemOptions(itemID, tenantID string) ([]map[string]any, error)
	CreateItemOption(itemID, tenantID string, body map[string]any) (map[string]any, error)
	ListOptionValues(optionID, tenantID string) ([]map[string]any, error)
	CreateOptionValue(optionID, tenantID string, body map[string]any) (map[string]any, error)

	// Tables
	GenerateTableQR(tableID, tenantID string) (string, error)
}

type adminMenuService struct {
	cat CategoryRepository
	it  ItemRepository
	op  OptionRepository
}

func NewAdminMenuService(cat CategoryRepository, it ItemRepository, op OptionRepository) AdminMenuService {
	return &adminMenuService{cat: cat, it: it, op: op}
}

// ---- Categories
func (s *adminMenuService) ListCategories(tenantID string) ([]map[string]any, error) {
	xs, err := s.cat.List(tenantID); if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(xs))
	for _, c := range xs {
		out = append(out, map[string]any{
			"id": c.ID, "tenant_id": c.TenantID, "name": c.Name, "sort": c.Sort, "is_active": c.IsActive,
		})
	}
	return out, nil
}
func (s *adminMenuService) CreateCategory(tenantID string, body map[string]any) (map[string]any, error) {
	c := domain.Category{
		TenantID: tenantID,
		Name: body["name"].(string),
	}
	if v, ok := body["sort"].(float64); ok { c.Sort = int(v) }
	if v, ok := body["is_active"].(bool); ok { c.IsActive = v }
	if err := s.cat.Create(&c); err != nil { return nil, err }
	return map[string]any{"id": c.ID, "tenant_id": c.TenantID, "name": c.Name, "sort": c.Sort, "is_active": c.IsActive}, nil
}
func (s *adminMenuService) ReplaceCategory(tenantID, id string, body map[string]any) (map[string]any, error) {
	c := domain.Category{
		ID: id, TenantID: tenantID, Name: body["name"].(string),
	}
	if v, ok := body["sort"].(float64); ok { c.Sort = int(v) }
	if v, ok := body["is_active"].(bool); ok { c.IsActive = v }
	if err := s.cat.Replace(&c); err != nil { return nil, err }
	return map[string]any{"id": c.ID, "tenant_id": c.TenantID, "name": c.Name, "sort": c.Sort, "is_active": c.IsActive}, nil
}
func (s *adminMenuService) PatchCategory(tenantID, id string, body map[string]any) (map[string]any, error) {
	obj, err := s.cat.Patch(tenantID, id, body); if err != nil { return nil, err }
	return map[string]any{"id": obj.ID, "tenant_id": obj.TenantID, "name": obj.Name, "sort": obj.Sort, "is_active": obj.IsActive}, nil
}
func (s *adminMenuService) DeleteCategory(tenantID, id string) error {
	return s.cat.Delete(tenantID, id)
}

// ---- Items
func (s *adminMenuService) ListItems(tenantID, categoryID string) ([]map[string]any, error) {
	xs, err := s.it.List(tenantID, categoryID); if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(xs))
	for _, m := range xs {
		out = append(out, map[string]any{
			"id": m.ID, "tenant_id": m.TenantID, "category_id": m.CategoryID, "name": m.Name,
			"description": m.Description, "price": m.Price, "photo_url": m.PhotoURL,
			"flags": m.Flags, "is_active": m.IsActive,
		})
	}
	return out, nil
}
func (s *adminMenuService) CreateItem(tenantID string, body map[string]any) (map[string]any, error) {
	i := domain.Item{
		TenantID: tenantID,
		CategoryID: body["category_id"].(string),
		Name: body["name"].(string),
	}
	if v, ok := body["description"].(string); ok { i.Description = &v }
	if v, ok := body["price"].(float64); ok { i.Price = int64(v) }
	if v, ok := body["photo_url"].(string); ok { i.PhotoURL = &v }
	if v, ok := body["is_active"].(bool); ok { i.IsActive = v }
	if err := s.it.Create(&i); err != nil { return nil, err }
	return map[string]any{
		"id": i.ID, "tenant_id": i.TenantID, "category_id": i.CategoryID, "name": i.Name,
		"description": i.Description, "price": i.Price, "photo_url": i.PhotoURL, "flags": i.Flags, "is_active": i.IsActive,
	}, nil
}
func (s *adminMenuService) ReplaceItem(tenantID, id string, body map[string]any) (map[string]any, error) {
	i := domain.Item{
		ID: id, TenantID: tenantID,
	}
	if v, ok := body["category_id"].(string); ok { i.CategoryID = v }
	if v, ok := body["name"].(string); ok { i.Name = v }
	if v, ok := body["description"].(string); ok { i.Description = &v } else { i.Description = nil }
	if v, ok := body["price"].(float64); ok { i.Price = int64(v) }
	if v, ok := body["photo_url"].(string); ok { i.PhotoURL = &v } else { i.PhotoURL = nil }
	if v, ok := body["is_active"].(bool); ok { i.IsActive = v }
	if err := s.it.Replace(&i); err != nil { return nil, err }
	return s.objItem(tenantID, id)
}
func (s *adminMenuService) PatchItem(tenantID, id string, body map[string]any) (map[string]any, error) {
	obj, err := s.it.Patch(tenantID, id, body); if err != nil { return nil, err }
	return map[string]any{
		"id": obj.ID, "tenant_id": obj.TenantID, "category_id": obj.CategoryID, "name": obj.Name,
		"description": obj.Description, "price": obj.Price, "photo_url": obj.PhotoURL, "flags": obj.Flags, "is_active": obj.IsActive,
	}, nil
}
func (s *adminMenuService) DeleteItem(tenantID, id string) error { return s.it.Delete(tenantID, id) }
func (s *adminMenuService) ToggleOOS(tenantID, id string, isActive bool) (map[string]any, error) {
	obj, err := s.it.ToggleActive(tenantID, id, isActive); if err != nil { return nil, err }
	return map[string]any{
		"id": obj.ID, "tenant_id": obj.TenantID, "category_id": obj.CategoryID, "name": obj.Name,
		"description": obj.Description, "price": obj.Price, "photo_url": obj.PhotoURL, "flags": obj.Flags, "is_active": obj.IsActive,
	}, nil
}
func (s *adminMenuService) objItem(tenantID, id string) (map[string]any, error) {
	obj, err := s.it.FindByID(tenantID, id); if err != nil { return nil, err }
	return map[string]any{
		"id": obj.ID, "tenant_id": obj.TenantID, "category_id": obj.CategoryID, "name": obj.Name,
		"description": obj.Description, "price": obj.Price, "photo_url": obj.PhotoURL, "flags": obj.Flags, "is_active": obj.IsActive,
	}, nil
}

// ---- Options
func (s *adminMenuService) ListItemOptions(itemID, tenantID string) ([]map[string]any, error) {
	xs, err := s.op.ListItemOptions(itemID, tenantID); if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(xs))
	for _, o := range xs {
		out = append(out, map[string]any{"id": o.ID, "item_id": o.ItemID, "name": o.Name, "type": o.Type, "required": o.Required})
	}
	return out, nil
}
func (s *adminMenuService) CreateItemOption(itemID, tenantID string, body map[string]any) (map[string]any, error) {
	o := domain.ItemOption{Name: body["name"].(string), Type: body["type"].(string)}
	if v, ok := body["required"].(bool); ok { o.Required = v }
	if err := s.op.CreateItemOption(itemID, tenantID, &o); err != nil { return nil, err }
	return map[string]any{"id": o.ID, "item_id": o.ItemID, "name": o.Name, "type": o.Type, "required": o.Required}, nil
}
func (s *adminMenuService) ListOptionValues(optionID, tenantID string) ([]map[string]any, error) {
	xs, err := s.op.ListOptionValues(optionID, tenantID); if err != nil { return nil, err }
	out := make([]map[string]any, 0, len(xs))
	for _, v := range xs {
		out = append(out, map[string]any{"id": v.ID, "option_id": v.OptionID, "label": v.Label, "delta_price": v.DeltaPrice})
	}
	return out, nil
}
func (s *adminMenuService) CreateOptionValue(optionID, tenantID string, body map[string]any) (map[string]any, error) {
	v := domain.ItemOptionValue{ Label: body["label"].(string) }
	if dp, ok := body["delta_price"].(float64); ok { v.DeltaPrice = int64(dp) }
	if err := s.op.CreateOptionValue(optionID, tenantID, &v); err != nil { return nil, err }
	return map[string]any{"id": v.ID, "option_id": v.OptionID, "label": v.Label, "delta_price": v.DeltaPrice}, nil
}

// ---- Tables
func (s *adminMenuService) GenerateTableQR(tableID, tenantID string) (string, error) {
	// TODO: implement generate URL (e.g., presigned S3) – sementara stub
	return "https://example.com/qr/"+tableID, nil
}
