package handler

import (
	"time"

	"qrmenu/internal/platform/logging"
	"qrmenu/internal/usecase"

	"github.com/go-playground/validator/v10" // <-- pastikan di go.mod
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	uc     *usecase.AuthUC
	v      *validator.Validate
	isProd bool
}

func NewAuthHandler(uc *usecase.AuthUC, isProd bool) *AuthHandler {
	return &AuthHandler{
		uc:     uc,
		v:      validator.New(), // <-- INI PENTING, biar nggak nil
		isProd: isProd,
	}
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		logging.HandlerError(c, "Auth.Login", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err)
		return fiber.ErrBadRequest
	}
	// kalau mau sementara tanpa validator, komentar 2 baris ini
	if err := h.v.Struct(req); err != nil {
		logging.HandlerError(c, "Auth.Login", "invalid payload", fiber.StatusBadRequest, "validation_failed", err, "email", req.Email)
		return fiber.ErrBadRequest
	}

	// assume AuthUC.Login returns (token string, tenantID string, error)
	token, tenantID, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		logging.HandlerError(c, "Auth.Login", "authentication failed", fiber.StatusUnauthorized, "invalid_credentials", err, "email", req.Email)
		return fiber.ErrUnauthorized
	}
	logging.HandlerInfo(c, "Auth.Login", "login successful", fiber.StatusOK, "auth_success", "email", req.Email, "tenant_id", tenantID)

	// persist token in an HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   h.isProd, // only secure in production
		SameSite: "Lax",    // safe default for local development
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // can be aligned with the JWT expiry
	})

	return c.JSON(fiber.Map{"ok": true})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   h.isProd,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
	logging.HandlerInfo(c, "Auth.Logout", "logout successful", fiber.StatusOK, "logout_success")
	return c.JSON(fiber.Map{"ok": true})
}
