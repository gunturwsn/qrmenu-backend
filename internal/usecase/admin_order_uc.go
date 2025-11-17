package usecase

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/repository"
)

type AdminOrdersUC struct {
	orders repository.OrderRepository
}

func NewAdminOrdersUC(r repository.OrderRepository) *AdminOrdersUC {
	return &AdminOrdersUC{orders: r}
}

func (u *AdminOrdersUC) List(tenantID, status, cursor string) (repository.OrdersPage, error) {
	logging.UsecaseInfo("AdminOrders.List", "listing orders", "orders_list_requested", "tenant_id", tenantID, "status", status, "cursor", cursor)
	page, err := u.orders.ListAdmin(tenantID, status, cursor, 20)
	if err != nil {
		logging.UsecaseError("AdminOrders.List", "repository error", "orders_list_failed", err, "tenant_id", tenantID, "status", status, "cursor", cursor)
		return repository.OrdersPage{}, err
	}
	logging.UsecaseInfo("AdminOrders.List", "orders loaded", "orders_listed", "tenant_id", tenantID, "status", status, "count", len(page.Data))
	return page, nil
}

func (u *AdminOrdersUC) UpdateStatus(tenantID, id, status string) (*domain.Order, error) {
	logging.UsecaseInfo("AdminOrders.UpdateStatus", "updating status", "order_status_update_requested", "tenant_id", tenantID, "order_id", id, "status", status)
	ord, err := u.orders.UpdateStatus(tenantID, id, domain.OrderStatus(status))
	if err != nil {
		logging.UsecaseError("AdminOrders.UpdateStatus", "repository error", "order_status_update_failed", err, "tenant_id", tenantID, "order_id", id, "status", status)
		return nil, err
	}
	logging.UsecaseInfo("AdminOrders.UpdateStatus", "status updated", "order_status_updated", "tenant_id", tenantID, "order_id", id, "status", status)
	return ord, nil
}
