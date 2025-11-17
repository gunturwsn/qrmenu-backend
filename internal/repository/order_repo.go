package repository

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
)

type OrdersPage struct {
	Data       []domain.Order `json:"data"`
	NextCursor *string        `json:"next_cursor"`
}

type OrderRepository interface {
	CreateGuestOrder(req domain.OrderCreateRequest) (orderID string, status domain.OrderStatus, err error)
	ListAdmin(tenantID, status, cursor string, limit int) (OrdersPage, error)
	UpdateStatus(tenantID, id string, status domain.OrderStatus) (*domain.Order, error)
}

type orderRepo struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) OrderRepository { return &orderRepo{db: db} }

// CreateGuestOrder creates an order from public endpoint payload.
// - Resolves tenant by code, table by token (must belong to tenant).
// - Creates order with WAITING & UNPAID status, then inserts items.
func (r *orderRepo) CreateGuestOrder(req domain.OrderCreateRequest) (string, domain.OrderStatus, error) {
	var orderID string
	logging.RepoInfo("OrderRepository.CreateGuestOrder", "create guest order", "order_create_requested", "tenant", req.Tenant, "table_token", req.TableToken, "items", len(req.Items))
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Resolve tenant by code
		var tenant domain.Tenant
		if err := tx.Where("code = ?", req.Tenant).First(&tenant).Error; err != nil {
			logging.RepoError("OrderRepository.CreateGuestOrder", "tenant lookup failed", "tenant_lookup_failed", err, "tenant_code", req.Tenant)
			return err
		}
		// Resolve table by token within tenant
		var table domain.Table
		if err := tx.Where("token = ? AND tenant_id = ? AND is_active = TRUE", req.TableToken, tenant.ID).
			First(&table).Error; err != nil {
			logging.RepoError("OrderRepository.CreateGuestOrder", "table lookup failed", "table_lookup_failed", err, "tenant_id", tenant.ID, "table_token", req.TableToken)
			return err
		}

		// Create order
		order := domain.Order{
			TenantID:     tenant.ID,
			TableID:      table.ID,
			GuestSession: req.GuestSession,
			Note:         req.Note,
			Status:       domain.OrderWaiting,
			PaidStatus:   domain.Unpaid,
		}
		if err := tx.Create(&order).Error; err != nil {
			logging.RepoError("OrderRepository.CreateGuestOrder", "order insert failed", "order_insert_failed", err, "tenant_id", tenant.ID)
			return err
		}
		orderID = order.ID

		// Create items
		for _, it := range req.Items {
			var menuItem domain.Item
			if err := tx.Where("id = ? AND tenant_id = ? AND is_active = TRUE", it.ItemID, tenant.ID).
				First(&menuItem).Error; err != nil {
				logging.RepoError("OrderRepository.CreateGuestOrder", "menu item lookup failed", "menu_item_lookup_failed", err, "tenant_id", tenant.ID, "item_id", it.ItemID)
				return err
			}
			oi := domain.OrderItem{
				OrderID:   order.ID,
				ItemID:    menuItem.ID,
				Name:      menuItem.Name,
				Qty:       it.Qty,
				UnitPrice: menuItem.Price,
			}
			if it.Options != nil {
				oi.Options = datatypes.JSONMap(it.Options) // jsonb
			}
			if err := tx.Create(&oi).Error; err != nil {
				logging.RepoError("OrderRepository.CreateGuestOrder", "order item insert failed", "order_item_insert_failed", err, "order_id", order.ID, "item_id", menuItem.ID)
				return err
			}
		}

		return nil
	})
	if err != nil {
		logging.RepoError("OrderRepository.CreateGuestOrder", "transaction failed", "transaction_failed", err, "tenant", req.Tenant)
		return orderID, domain.OrderWaiting, err
	}
	logging.RepoInfo("OrderRepository.CreateGuestOrder", "order created", "order_created", "order_id", orderID, "tenant", req.Tenant)
	return orderID, domain.OrderWaiting, nil
}

// ListAdmin returns paginated orders for a tenant with optional status filter.
// Cursor format: base64url("RFC3339Nano|order_id") and paging by (created_at,id) DESC.
func (r *orderRepo) ListAdmin(tenantID, status, cursor string, limit int) (OrdersPage, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.Where("tenant_id = ?", tenantID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	// Apply cursor: fetch records older than the cursor (DESC ordering)
	if cursor != "" {
		if ts, id, ok := decodeCursor(cursor); ok {
			q = q.Where("(created_at, id) < (?, ?)", ts, id)
		}
	}

	var rows []domain.Order
	if err := q.Order("created_at DESC, id DESC").
		Limit(limit + 1).
		Preload("Items").
		Find(&rows).Error; err != nil {
		logging.RepoError("OrderRepository.ListAdmin", "query failed", "query_failed", err, "tenant_id", tenantID, "status", status)
		return OrdersPage{}, err
	}

	var next *string
	if len(rows) > limit {
		last := rows[limit-1]
		cur := encodeCursor(last.CreatedAt, last.ID)
		next = &cur
		rows = rows[:limit]
	}

	logging.RepoInfo("OrderRepository.ListAdmin", "orders listed", "orders_listed", "tenant_id", tenantID, "status", status, "count", len(rows), "has_next", next != nil)
	return OrdersPage{Data: rows, NextCursor: next}, nil
}

// UpdateStatus updates an order status (scoped by tenant) and returns the updated order.
func (r *orderRepo) UpdateStatus(tenantID, id string, status domain.OrderStatus) (*domain.Order, error) {
	if err := r.db.Model(&domain.Order{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("status", status).Error; err != nil {
		logging.RepoError("OrderRepository.UpdateStatus", "update failed", "update_failed", err, "tenant_id", tenantID, "order_id", id, "status", status)
		return nil, err
	}
	var o domain.Order
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Items").
		First(&o).Error; err != nil {
		logging.RepoError("OrderRepository.UpdateStatus", "load failed", "load_failed", err, "tenant_id", tenantID, "order_id", id)
		return nil, err
	}
	logging.RepoInfo("OrderRepository.UpdateStatus", "order updated", "order_updated", "tenant_id", tenantID, "order_id", id, "status", o.Status)
	return &o, nil
}

// --- cursor helpers ---

func encodeCursor(t time.Time, id string) string {
	raw := fmt.Sprintf("%s|%s", t.UTC().Format(time.RFC3339Nano), id)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(s string) (time.Time, string, bool) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return time.Time{}, "", false
	}
	parts := strings.SplitN(string(b), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, "", false
	}
	ts, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", false
	}
	return ts, parts[1], true
}
