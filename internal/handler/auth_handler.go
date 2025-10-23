package handler

import (
	"time"

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
		return fiber.ErrBadRequest
	}
	// kalau mau sementara tanpa validator, komentar 2 baris ini
	if err := h.v.Struct(req); err != nil {
		return fiber.ErrBadRequest
	}

	// asumsi AuthUC.Login return (token string, err error)
	token, _, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	// set HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   h.isProd,          // hanya secure di prod
		SameSite: "Lax",             // aman untuk local dev
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // atau pakai TTL JWT kalau kamu expose
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
	return c.JSON(fiber.Map{"ok": true})
}
