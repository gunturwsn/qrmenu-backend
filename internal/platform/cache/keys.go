package cache

import "fmt"

func KeyMenusByTenant(tenantID string) string {
	return fmt.Sprintf("menus:tenant:%s", tenantID)
}

func KeyMenuByID(menuID string) string {
	return fmt.Sprintf("menu:%s", menuID)
}
