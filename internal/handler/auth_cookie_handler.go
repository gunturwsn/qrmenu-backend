package handler

import (
	"qrmenu/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthCookieHandler struct {
	uc     *usecase.AuthUC
	isProd bool
}

func NewAuthCookieHandler(uc *usecase.AuthUC, isProd bool) *AuthCookieHandler {
	return &AuthCookieHandler{uc: uc, isProd: isProd}
}

func (h *AuthCookieHandler) Login(c *fiber.Ctx) error {
	var in struct{ Email, Password string }
	if err := c.BodyParser(&in); err != nil { return fiber.ErrBadRequest }
	jwt, _, err := h.uc.Login(in.Email, in.Password)
	if err != nil { return fiber.ErrUnauthorized }
	c.Cookie(&fiber.Cookie{
		Name: "admin_token", Value: jwt, Path: "/", HTTPOnly: true,
		SameSite: "Lax", Secure: h.isProd, Expires: time.Now().Add(24*time.Hour),
	})
	return c.JSON(fiber.Map{"ok": true})
}
func (h *AuthCookieHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{Name:"admin_token",Value:"",Path:"/",HTTPOnly:true,Expires:time.Unix(0,0)})
	return c.JSON(fiber.Map{"ok": true})
}
