package handler

import (
	"time"

	"qrmenu/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SetupHandler struct {
	uc *usecase.SetupUC
	v  *validator.Validate
}

func NewSetupHandler(uc *usecase.SetupUC) *SetupHandler {
	return &SetupHandler{uc: uc, v: validator.New()}
}

type setupTenantReq struct {
	TenantCode string `json:"tenant_code" validate:"required"`
	TenantName string `json:"tenant_name"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
}

type setupStatusResp struct {
	Scope      string `json:"scope"`
	Initialized bool  `json:"initialized"`
	Tenant     string `json:"tenant,omitempty"`
}

// GET /setup/status?tenant=code
func (h *SetupHandler) Status(c *fiber.Ctx) error {
	code := c.Query("tenant_code")
	if code == "" {
		// global view tidak relevan untuk model B → kembalikan false saja atau info
		return c.JSON(setupStatusResp{Scope: "tenant", Initialized: false})
	}
	init, err := h.uc.IsTenantInitialized(code)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(setupStatusResp{Scope: "tenant", Initialized: init, Tenant: code})
}

// POST /setup/admin (per-tenant)
func (h *SetupHandler) SetupTenant(c *fiber.Ctx) error {
	// Aktifkan guard token jika mau
	// setupToken := c.Get("X-Setup-Token")
	// if setupToken == "" || setupToken != os.Getenv("SETUP_TOKEN") { return fiber.ErrForbidden }

	var req setupTenantReq
	if err := c.BodyParser(&req); err != nil { return fiber.ErrBadRequest }
	if err := h.v.Struct(req); err != nil { return fiber.ErrBadRequest }

	token, err := h.uc.SetupFirstAdminForTenant(usecase.SetupTenantRequest{
		TenantCode: req.TenantCode,
		TenantName: req.TenantName,
		Email:      req.Email,
		Password:   req.Password,
	})
	if err != nil {
		if err.Error() == "tenant already initialized" {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// set cookie admin_token supaya langsung bisa akses /admin/* untuk tenant tsb
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true})
}