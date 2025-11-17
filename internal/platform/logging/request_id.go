package logging

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// RequestIDFromFiber extracts the request identifier injected by the requestid middleware.
// Falls back to "-" when unavailable so callers can still include a placeholder in logs.
func RequestIDFromFiber(c *fiber.Ctx) string {
	if c == nil {
		return "-"
	}
	if id := c.Get(fiber.HeaderXRequestID); id != "" {
		return id
	}
	if id := c.GetRespHeader(fiber.HeaderXRequestID); id != "" {
		return id
	}
	if v := c.Locals(requestid.ConfigDefault.ContextKey); v != nil {
		if id, ok := v.(string); ok && id != "" {
			return id
		}
	}
	return "-"
}
