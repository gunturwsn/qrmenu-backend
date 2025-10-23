package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminCookieOnly(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Cookies("admin_token")
		if raw == "" {
			return fiber.ErrUnauthorized
		}
		tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			return fiber.ErrUnauthorized
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok || claims["role"] != "admin" {
			return fiber.ErrForbidden
		}
		c.Locals("admin_id", claims["sub"])
		c.Locals("admin_email", claims["email"])
		c.Locals("tenant_id", claims["tenant_id"])
		return c.Next()
	}
}
