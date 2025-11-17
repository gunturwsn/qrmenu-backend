package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	redisstore "github.com/gofiber/storage/redis/v3"

	"qrmenu/internal/config"
	"qrmenu/internal/platform/conv"
	"qrmenu/internal/platform/logging"
)

// Fungsi utama dipanggil dari main.go
func RegisterHTTP(app *fiber.App, extraOrigins []string) {
	cfg := config.Load()
	host, port := conv.SplitHostPort(cfg.Redis.Addr, "redis", 6379)

	// Call the retry helper that keeps dialing Redis until it succeeds.
	store := newRedisStoreWithRetry(host, port, cfg.Redis.Password, cfg.Redis.DB)
	_ = store // reserved for rate-limiter/session usage if needed later.

	log.Println("[redis] connection established to", host, port)

	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true, StackTraceHandler: logStackTrace}))
	app.Use(httpAccessLogger())
}

// Fungsi helper retry koneksi Redis
func newRedisStoreWithRetry(host string, port int, pass string, db int) *redisstore.Storage {
	cfg := redisstore.Config{
		Host: host, Port: port, Password: pass, Database: db,
	}
	var lastErr error
	for i := 0; i < 10; i++ { // attempt up to 10 retries
		st := redisstore.New(cfg)
		if err := pingRedis(st); err == nil {
			return st
		} else {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
		}
	}
	panic(lastErr) // follow-up option: replace with log.Fatal(lastErr)
}

// pingRedis performs a lightweight read/write to verify Redis connectivity.
func pingRedis(store *redisstore.Storage) error {
	const key = "health:ping"
	if err := store.Set(key, []byte("ok"), 5*time.Second); err != nil {
		return err
	}
	_, err := store.Get(key)
	return err
}

func httpAccessLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		if status == 0 {
			// Fiber keeps status 0 for some errors until send; normalize for logging.
			if err != nil {
				status = fiber.StatusInternalServerError
			} else {
				status = fiber.StatusOK
			}
		}

		reqID := logging.RequestIDFromFiber(c)
		latency := time.Since(start)

		log.Printf("[http] request_id=%s method=%s path=%s status=%d latency=%s ip=%s",
			reqID, c.Method(), c.OriginalURL(), status, latency, c.IP())

		if err != nil {
			log.Printf("[http] request_id=%s error=%v", reqID, err)
		}

		if status >= fiber.StatusBadRequest {
			if body := trimLogBody(c.Body()); body != "" {
				log.Printf("[http] request_id=%s request_body=%s", reqID, body)
			}
			if body := trimLogBody(c.Response().Body()); body != "" {
				log.Printf("[http] request_id=%s response_body=%s", reqID, body)
			}
		}

		return err
	}
}

func trimLogBody(body []byte) string {
	const limit = 512
	if len(body) == 0 {
		return ""
	}
	text := string(body)
	text = strings.ReplaceAll(text, "\n", "\\n")
	if len(text) > limit {
		return text[:limit] + "...(truncated)"
	}
	return text
}

func logStackTrace(c *fiber.Ctx, e interface{}) {
	reqID := logging.RequestIDFromFiber(c)
	log.Printf("[panic] request_id=%s status_code=500 error_code=panic error=%v", reqID, e)
}
