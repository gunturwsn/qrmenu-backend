package usecase

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"qrmenu/internal/domain"
	"qrmenu/internal/platform/cache"
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
		return nil, errors.New("invalid tenant code")
	}

	key := cache.KeyMenusByTenant(code)

	if u.cache != nil && u.ttl > 0 {
		if cached, err := u.cache.Get(key); err == nil && cached != "" {
			var resp domain.MenuResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				return &resp, nil
			}
		}
	}

	menu, err := u.query.GetMenuByTenantCode(code)
	if err != nil {
		return nil, err
	}

	if u.cache != nil && u.ttl > 0 {
		if payload, err := json.Marshal(menu); err == nil {
			_ = u.cache.Set(key, string(payload), u.ttl)
		}
	}

	return menu, nil
}

func (u *menuUC) InvalidateTenantMenu(code string) {
	code = strings.TrimSpace(code)
	if code == "" || u.cache == nil {
		return
	}
	_ = u.cache.Del(cache.KeyMenusByTenant(code))
}
