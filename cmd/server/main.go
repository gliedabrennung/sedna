package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gliedabrennung/messenger-core/internal/config"
	"github.com/gliedabrennung/messenger-core/internal/controller/http"
	"github.com/gliedabrennung/messenger-core/internal/messenger"
	"github.com/gliedabrennung/messenger-core/internal/repository/postgres"
	"github.com/gliedabrennung/messenger-core/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbpool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer dbpool.Close()

	if sqlBytes, err := os.ReadFile("migrations/0001.sql"); err == nil {
		if _, err := dbpool.Exec(ctx, string(sqlBytes)); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	} else {
		log.Printf("WARNING: could not read migration file: %v", err)
	}

	repo := postgres.NewPostgresRepository(dbpool)
	authUseCase := usecase.NewAuthUseCase(repo, cfg.JWTSecret, cfg.JWTTTL)

	hub := messenger.NewHub()
	go hub.Run(ctx)

	allowedOrigins := []string{"*"} // TODO: configure per environment
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		allowedOrigins = strings.Split(v, ",")
	}

	h := server.Default(
		server.WithHostPorts(cfg.Addr),
		server.WithHandleMethodNotAllowed(true),
	)

	http.SetupRouter(h, http.Deps{
		Auth:           authUseCase,
		Hub:            hub,
		JWTSecret:      cfg.JWTSecret,
		AllowedOrigins: allowedOrigins,
	})

	h.Spin()
}
