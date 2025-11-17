package handler

import (
	"time"

	"qrmenu/internal/platform/logging"
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
	Scope       string `json:"scope"`
	Initialized bool   `json:"initialized"`
	Tenant      string `json:"tenant,omitempty"`
}

// GET /setup/status?tenant=code
func (h *SetupHandler) Status(c *fiber.Ctx) error {
	code := c.Query("tenant_code")
	if code == "" {
		// A tenant-agnostic status check is not applicable; return a default response.
		logging.HandlerInfo(c, "Setup.Status", "tenant code missing, returning default", fiber.StatusOK, "status_default")
		return c.JSON(setupStatusResp{Scope: "tenant", Initialized: false})
	}
	init, err := h.uc.IsTenantInitialized(code)
	if err != nil {
		logging.HandlerError(c, "Setup.Status", "failed to check initialization", fiber.StatusBadRequest, "status_lookup_failed", err, "tenant_code", code)
		return fiber.ErrBadRequest
	}
	logging.HandlerInfo(c, "Setup.Status", "status fetched", fiber.StatusOK, "status_success", "tenant_code", code, "initialized", init)
	return c.JSON(setupStatusResp{Scope: "tenant", Initialized: init, Tenant: code})
}

// POST /setup/admin (per-tenant)
func (h *SetupHandler) SetupTenant(c *fiber.Ctx) error {
	// Enable a guard token here if you want to protect the setup endpoint:
	// setupToken := c.Get("X-Setup-Token")
	// if setupToken == "" || setupToken != os.Getenv("SETUP_TOKEN") { return fiber.ErrForbidden }

	var req setupTenantReq
	if err := c.BodyParser(&req); err != nil {
		logging.HandlerError(c, "Setup.SetupTenant", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err)
		return fiber.ErrBadRequest
	}
	if err := h.v.Struct(req); err != nil {
		logging.HandlerError(c, "Setup.SetupTenant", "validation failed", fiber.StatusBadRequest, "validation_failed", err, "tenant_code", req.TenantCode, "email", req.Email)
		return fiber.ErrBadRequest
	}

	token, err := h.uc.SetupFirstAdminForTenant(usecase.SetupTenantRequest{
		TenantCode: req.TenantCode,
		TenantName: req.TenantName,
		Email:      req.Email,
		Password:   req.Password,
	})
	if err != nil {
		if err.Error() == "tenant already initialized" {
			logging.HandlerError(c, "Setup.SetupTenant", "tenant already initialized", fiber.StatusForbidden, "tenant_already_initialized", err, "tenant_code", req.TenantCode)
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		logging.HandlerError(c, "Setup.SetupTenant", "failed to setup tenant", fiber.StatusBadRequest, "setup_failed", err, "tenant_code", req.TenantCode)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Set admin_token so the newly created admin can access /admin/* immediately.
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
	logging.HandlerInfo(c, "Setup.SetupTenant", "tenant initialized", fiber.StatusCreated, "tenant_initialized", "tenant_code", req.TenantCode, "email", req.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true})
}
