package usecase

import (
	"qrmenu/internal/repository"
)

type TableUC struct {
	tableRepo repository.TableRepository
}

func NewTableUC(r repository.TableRepository) *TableUC { return &TableUC{tableRepo: r} }

func (u *TableUC) ResolveByToken(token string) (tenant any, table any, err error) {
	t, tb, e := u.tableRepo.ResolveByToken(token)
	return t, tb, e
}
