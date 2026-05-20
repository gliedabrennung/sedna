package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	"github.com/gliedabrennung/messenger-core/internal/config"
	"github.com/gliedabrennung/messenger-core/internal/controller/http"
	"github.com/gliedabrennung/messenger-core/internal/controller/http/middleware"
	"github.com/gliedabrennung/messenger-core/internal/messenger"
	"github.com/gliedabrennung/messenger-core/internal/repository/postgres"
	"github.com/gliedabrennung/messenger-core/internal/usecase"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var addr = ":8080"

// @title Messenger Core API
// @version 1.0
// @description This is a REST API for the messenger application.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()

	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		hlog.Errorf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	m, err := migrate.New(
		"file://migrations",
		cfg.DSN,
	)
	if err != nil {
		hlog.Errorf("Unable to create migrate instance: %v\n", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		hlog.Errorf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	repo := postgres.NewPostgresRepository(dbpool)
	authUseCase := usecase.NewAuthUseCase(repo, cfg.JWTSecret, cfg.JWTTTL)
	msgUseCase := usecase.NewMessageUseCase(repo)
	contactUseCase := usecase.NewContactUseCase(repo)

	messenger.InitHub(msgUseCase)

	os.MkdirAll("./uploads", 0755)

	h := server.Default(
		server.WithHostPorts(addr),
		server.WithHandleMethodNotAllowed(true),
	)

	h.Use(middleware.RequestLogger())
	h.Use(middleware.RateLimiter())

	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h.Static("/uploads", "./")

	http.SetupRouter(h, authUseCase, msgUseCase, contactUseCase, cfg.JWTSecret)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := h.Run(); err != nil {
			hlog.Errorf("Hertz server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	hlog.Info("Shutting down server gracefully...")

	// Timeout for shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Shutdown(shutdownCtx); err != nil {
		hlog.Errorf("Hertz graceful shutdown failed: %v", err)
	}

	hlog.Info("Hertz server stopped. Closing DB connection pool...")
}
