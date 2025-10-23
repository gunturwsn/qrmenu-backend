package repository

import "qrmenu/internal/domain"

type AdminOrdersQuery interface {
	List(status, cursor, tenantID string) (OrdersPage, error)
	UpdateStatus(id, status, tenantID string) (map[string]any, error)
}

type adminOrdersQuery struct {
	orders OrderRepository
}

func NewAdminOrdersQuery(orders OrderRepository) AdminOrdersQuery {
	return &adminOrdersQuery{orders: orders}
}

func (q *adminOrdersQuery) List(status, cursor, tenantID string) (OrdersPage, error) {
	return q.orders.ListAdmin(tenantID, status, cursor, 20)
}

func (q *adminOrdersQuery) UpdateStatus(id, status, tenantID string) (map[string]any, error) {
	o, err := q.orders.UpdateStatus(tenantID, id, domain.OrderStatus(status))
	if err != nil { return map[string]any{}, err }
	return map[string]any{
		"id": o.ID, "tenant_id": o.TenantID, "table_id": o.TableID,
		"status": o.Status, "paid_status": o.PaidStatus, "note": o.Note,
		"created_at": o.CreatedAt, "items": o.Items,
	}, nil
}
