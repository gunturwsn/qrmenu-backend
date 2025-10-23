package usecase

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/repository"
)

type AdminOrdersUC struct {
	orders repository.OrderRepository
}

func NewAdminOrdersUC(r repository.OrderRepository) *AdminOrdersUC {
	return &AdminOrdersUC{orders: r}
}

func (u *AdminOrdersUC) List(tenantID, status, cursor string) (repository.OrdersPage, error) {
	return u.orders.ListAdmin(tenantID, status, cursor, 20)
}


func (u *AdminOrdersUC) UpdateStatus(tenantID, id, status string) (*domain.Order, error) {
	return u.orders.UpdateStatus(tenantID, id, domain.OrderStatus(status))
}
