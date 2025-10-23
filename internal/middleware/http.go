package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	redisstore "github.com/gofiber/storage/redis/v3"

	"qrmenu/internal/config"
	"qrmenu/internal/platform/conv"
)

// Fungsi utama dipanggil dari main.go
func RegisterHTTP(app *fiber.App, extraOrigins []string) {
	cfg := config.Load()
	host, port := conv.SplitHostPort(cfg.Redis.Addr, "redis", 6379)

	// 🔥 Panggil helper retry di sini
	store := newRedisStoreWithRetry(host, port, cfg.Redis.Password, cfg.Redis.DB)
	_ = store // contoh, bisa kamu gunakan untuk limiter/session dsb.

	log.Println("[redis] connection established to", host, port)

	// ... lanjutkan middleware kamu seperti CORS, logger, recover, dsb.
}

// Fungsi helper retry koneksi Redis
func newRedisStoreWithRetry(host string, port int, pass string, db int) *redisstore.Storage {
	cfg := redisstore.Config{
		Host: host, Port: port, Password: pass, Database: db,
	}
	var lastErr error
	for i := 0; i < 10; i++ { // coba sampai 10x
		st := redisstore.New(cfg)
		if err := pingRedis(st); err == nil {
			return st
		} else {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
		}
	}
	panic(lastErr) // atau log.Fatal(lastErr)
}

// Ping sederhana untuk test koneksi Redis
func pingRedis(store *redisstore.Storage) error {
	const key = "health:ping"
	if err := store.Set(key, []byte("ok"), 5*time.Second); err != nil {
		return err
	}
	_, err := store.Get(key)
	return err
}
