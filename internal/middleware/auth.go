package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminOnly(jwtSecret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        raw := c.Cookies("admin_token")
        if raw == "" { return fiber.ErrUnauthorized }

        tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.ErrUnauthorized
            }
            return []byte(jwtSecret), nil
        })
        if err != nil || !tok.Valid { return fiber.ErrUnauthorized }

        claims, ok := tok.Claims.(jwt.MapClaims)
        if !ok { return fiber.ErrUnauthorized }

        // role wajib admin
        if role, _ := claims["role"].(string); role != "admin" {
            return fiber.ErrForbidden
        }

        // tenant_id wajib ada
        tenantID, _ := claims["tenant_id"].(string)
        if tenantID == "" { 
            return fiber.ErrForbidden
        }

        // set ke context untuk handler
        if sub, _ := claims["sub"].(string); sub != "" {
            c.Locals("admin_id", sub)
        }
        c.Locals("tenant_id", tenantID)
        if email, _ := claims["email"].(string); email != "" {
            c.Locals("email", email)
        }
        return c.Next()
    }
}


// Simple CORS config sudah ada di main; opsional cors.go terpisah kalau mau granular
func LogErr(err error) { if err != nil { log.Println(err) } }
