package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastlearn/backend/api"
	"github.com/fastlearn/backend/api/handlers"
	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/db"
	"github.com/fastlearn/backend/logger"
	"github.com/fastlearn/backend/queue"
	"github.com/fastlearn/backend/repository"
	"github.com/fastlearn/backend/services"
	"github.com/joho/godotenv"
)

func main() {
	// Ищем .env в текущей папке, потом на уровень выше (запуск из backend/ или из корня)
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load("../.env")
	}

	logger.Init()

	cfg := config.Load()

	// ── Infrastructure ────────────────────────────────────────────
	database, err := db.Connect(cfg.DatabaseURL())
	if err != nil {
		slog.Error("db connect failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	//migrationsDir := cfg.MigrationsDir
	//if err := db.RunMigrations(context.Background(), database, migrationsDir); err != nil {
	//	slog.Error("migrations failed", "error", err)
	//	os.Exit(1)
	//}

	rdb, err := db.ConnectRedis(cfg.RedisURL())
	if err != nil {
		slog.Error("redis connect failed", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	// ── Repositories ──────────────────────────────────────────────
	userRepo := repository.NewUserRepo(database)
	videoRepo := repository.NewVideoRepo(database)
	subRepo := repository.NewSubscriptionRepo(database)
	summaryRepo := repository.NewSummaryRepo(database)
	procRepo := repository.NewProcessingRepo(cfg.FastAPIURL)

	// ── Services ──────────────────────────────────────────────────
	authSvc := services.NewAuthService(userRepo, subRepo, cfg)
	videoSvc := services.NewVideoService(videoRepo, summaryRepo, subRepo, rdb, cfg)
	subSvc := services.NewSubscriptionService(subRepo, cfg)

	// ── Handlers ──────────────────────────────────────────────────
	authH := handlers.NewAuthHandler(authSvc)
	videoH := handlers.NewVideoHandler(videoSvc)
	subH := handlers.NewSubscriptionHandler(subSvc)

	// ── Background worker ─────────────────────────────────────────
	processor := queue.NewProcessor(videoRepo, summaryRepo, procRepo, rdb, cfg)
	go processor.Start()

	// ── HTTP Server ───────────────────────────────────────────────
	e := api.NewRouter(authH, videoH, subH, cfg)

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := e.Start(":" + cfg.Port); err != nil {
			slog.Info("server stopped", "reason", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
