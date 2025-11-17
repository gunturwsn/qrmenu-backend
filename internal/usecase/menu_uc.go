package usecase

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"qrmenu/internal/domain"
	"qrmenu/internal/platform/cache"
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/repository"
)

type MenuUC interface {
	GetMenuByTenantCode(code string) (*domain.MenuResponse, error)
	InvalidateTenantMenu(code string)
}

type menuUC struct {
	query repository.MenuQuery
	cache cache.Cache
	ttl   time.Duration
}

func NewMenuUC(q repository.MenuQuery, rc cache.Cache, ttl time.Duration) MenuUC {
	return &menuUC{query: q, cache: rc, ttl: ttl}
}

func (u *menuUC) GetMenuByTenantCode(code string) (*domain.MenuResponse, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		err := errors.New("invalid tenant code")
		logging.UsecaseError("Menu.GetMenuByTenantCode", "tenant code empty", "invalid_tenant_code", err)
		return nil, err
	}

	key := cache.KeyMenusByTenant(code)

	if u.cache != nil && u.ttl > 0 {
		if cached, err := u.cache.Get(key); err == nil && cached != "" {
			var resp domain.MenuResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				logging.UsecaseInfo("Menu.GetMenuByTenantCode", "cache hit", "cache_hit", "tenant_code", code)
				return &resp, nil
			}
		}
	}

	menu, err := u.query.GetMenuByTenantCode(code)
	if err != nil {
		logging.UsecaseError("Menu.GetMenuByTenantCode", "query failed", "menu_query_failed", err, "tenant_code", code)
		return nil, err
	}

	if u.cache != nil && u.ttl > 0 {
		if payload, err := json.Marshal(menu); err == nil {
			if err := u.cache.Set(key, string(payload), u.ttl); err != nil {
				logging.UsecaseError("Menu.GetMenuByTenantCode", "cache set failed", "cache_set_failed", err, "tenant_code", code)
			}
		} else {
			logging.UsecaseError("Menu.GetMenuByTenantCode", "failed to marshal cache payload", "cache_marshal_failed", err, "tenant_code", code)
		}
	}

	logging.UsecaseInfo("Menu.GetMenuByTenantCode", "menu fetched", "menu_fetched", "tenant_code", code, "categories", len(menu.Categories), "items", len(menu.Items))
	return menu, nil
}

func (u *menuUC) InvalidateTenantMenu(code string) {
	code = strings.TrimSpace(code)
	if code == "" || u.cache == nil {
		return
	}
	if err := u.cache.Del(cache.KeyMenusByTenant(code)); err != nil {
		logging.UsecaseError("Menu.InvalidateTenantMenu", "cache delete failed", "cache_delete_failed", err, "tenant_code", code)
		return
	}
	logging.UsecaseInfo("Menu.InvalidateTenantMenu", "cache invalidated", "cache_invalidated", "tenant_code", code)
}
