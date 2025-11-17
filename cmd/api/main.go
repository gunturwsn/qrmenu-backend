package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"qrmenu/internal/config"
	"qrmenu/internal/handler"
	"qrmenu/internal/middleware"
	"qrmenu/internal/platform/cache"
	"qrmenu/internal/platform/db"
	"qrmenu/internal/platform/security"
	"qrmenu/internal/repository"
	"qrmenu/internal/transport/http"
	"qrmenu/internal/usecase"
)

func main() {
	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "run DB migrations startup")
	flag.Parse()

	// Load config
	cfg := config.Load()

	// Optimize scheduler to all available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Force runtime TZ → Asia/Jakarta
	if loc, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		time.Local = loc
	}

	// Connect DB
	gdb := db.Connect(cfg)

	// Auto migrate & seed admin
	if migrate {
		log.Println("DB migrations are managed externally via golang-migrate.")
		return
	}

	// Connect Redis
	rc := cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	defer func() {
		if err := rc.Close(); err != nil {
			log.Printf("warn: failed to close redis connection: %v", err)
		}
	}()
	defaultTTL := time.Duration(cfg.Redis.TTLSeconds) * time.Second

	// ===== Repositories =====
	adminRepo := repository.NewAdminRepository(gdb)
	tenantRepo := repository.NewTenantRepository(gdb)
	_ = tenantRepo
	tableRepo := repository.NewTableRepository(gdb)
	catRepo := repository.NewCategoryRepository(gdb)
	itemRepo := repository.NewItemRepository(gdb)
	optRepo := repository.NewOptionRepository(gdb)
	orderRepo := repository.NewOrderRepository(gdb)
	menuQuery := repository.NewMenuQuery(gdb)

	// ===== Security / JWT =====
	jwtMaker := security.NewJWT(cfg.JWTSecret, cfg.JWTExpiresMinute)

	// ===== Usecases =====
	setupUC := usecase.NewSetupUC(adminRepo, tenantRepo, jwtMaker)
	authUC := usecase.NewAuthUC(adminRepo, jwtMaker)
	menuUC := usecase.NewMenuUC(menuQuery, rc, defaultTTL)
	tableUC := usecase.NewTableUC(tableRepo)
	orderUC := usecase.NewOrderUC(orderRepo)
	adminMenuUC := usecase.NewAdminMenuUC(catRepo, itemRepo, optRepo)
	adminOrdersUC := usecase.NewAdminOrdersUC(orderRepo)

	// ===== Handlers =====
	setupH := handler.NewSetupHandler(setupUC)
	authH := handler.NewAuthHandler(authUC, cfg.IsProd())
	menuH := handler.NewMenuHandler(menuUC)
	tableH := handler.NewTableHandler(tableUC)
	orderPubH := handler.NewOrderPublicHandler(orderUC)

	// Admin menu handler now consumes the use case directly.
	adminMenuH := handler.NewAdminMenuHandler(adminMenuUC)

	// Admin orders handler can invoke the use case directly.
	adminOrdersH := handler.NewAdminOrdersHandler(adminOrdersUC)

	// ===== Fiber app =====
	app := fiber.New(fiber.Config{
		AppName:       cfg.AppName,
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Global middlewares (CORS, Recover, Logger)
	middleware.RegisterHTTP(app, cfg.AllowedOrigins)

	// ===== Routes =====
	http.Register(app, http.Deps{
		Auth:      authH,
		Menu:      menuH,
		Table:     tableH,
		OrderPub:  orderPubH,
		AdminMenu: adminMenuH,
		AdminOrd:  adminOrdersH,
		Setup:     setupH,
		JWTSecret: cfg.JWTSecret,
	})
	handler.RegisterSwaggerUI(app)

	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("failed to obtain sql.DB: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("warn: failed to close DB connection: %v", cerr)
		}
	}()

	// Graceful shutdown handling.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		log.Println("shutdown signal received, closing server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Printf("error shutting down server: %v", err)
		}
	}()

	// Start
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 %s running on %s (env=%s, tz=%s)", cfg.AppName, addr, cfg.AppEnv, time.Now().Location())
	if err := app.Listen(addr); err != nil {
		log.Printf("server closed: %v", err)
	}
}
