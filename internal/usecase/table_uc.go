package usecase

import (
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/repository"
)

type TableUC struct {
	tableRepo repository.TableRepository
}

func NewTableUC(r repository.TableRepository) *TableUC { return &TableUC{tableRepo: r} }

func (u *TableUC) ResolveByToken(token string) (tenant any, table any, err error) {
	logging.UsecaseInfo("Table.ResolveByToken", "resolving table", "table_resolve_requested", "token", token)
	t, tb, e := u.tableRepo.ResolveByToken(token)
	if e != nil {
		logging.UsecaseError("Table.ResolveByToken", "repository error", "table_resolve_failed", e, "token", token)
		return nil, nil, e
	}
	logging.UsecaseInfo("Table.ResolveByToken", "table resolved", "table_resolved", "token", token)
	return t, tb, nil
}
