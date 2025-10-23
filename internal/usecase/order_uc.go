package usecase

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/repository"
)

type OrderUC struct {
	repo repository.OrderRepository
}

func NewOrderUC(r repository.OrderRepository) *OrderUC { return &OrderUC{repo: r} }

// CreateGuestOrder meneruskan payload domain langsung ke repository.
// Return: orderID, status(string), error
func (u *OrderUC) CreateGuestOrder(req domain.OrderCreateRequest) (string, string, error) {
	id, st, err := u.repo.CreateGuestOrder(req)
	if err != nil {
		return "", "", err
	}
	return id, string(st), nil
}
