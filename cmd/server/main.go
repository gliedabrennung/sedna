package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gliedabrennung/sedna/internal/common/logger"
	"github.com/gliedabrennung/sedna/internal/config"
	"github.com/gliedabrennung/sedna/internal/controller/http"
	"github.com/gliedabrennung/sedna/internal/domain"
	"github.com/gliedabrennung/sedna/internal/repository/message"
	"github.com/gliedabrennung/sedna/internal/repository/postgres"
	"github.com/gliedabrennung/sedna/internal/usecase"
	"github.com/gliedabrennung/sedna/internal/ws"
	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		logger.Fatalf("startup failed: %v", err)
	}
}

func run() error {
	var (
		msgRepo       domain.MessageRepository
		scyllaSession *gocql.Session
	)

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return fmt.Errorf("create connection pool: %w", err)
	}
	defer dbpool.Close()

	if err := message.InitSchema(ctx, cfg.ScyllaHosts, cfg.ScyllaKeyspace); err != nil {
		logger.Warnf("warning: could not initialize scylla schema: %v", err)
	} else {
		cluster := gocql.NewCluster(cfg.ScyllaHosts...)
		cluster.Keyspace = cfg.ScyllaKeyspace
		cluster.Timeout = 5 * time.Second
		var err error
		scyllaSession, err = cluster.CreateSession()
		if err != nil {
			logger.Warnf("warning: could not connect to scylla (skipping messages feature): %v", err)
			scyllaSession = nil
		} else {
			defer scyllaSession.Close()
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warnf("warning: could not connect to redis: %v", err)
		rdb = nil
	} else {
		defer func() {
			if err := rdb.Close(); err != nil {
				logger.Warnf("warning: could not close redis connection: %v", err)
			}
		}()
	}

	if scyllaSession != nil && rdb != nil {
		msgRepo = message.NewRepository(scyllaSession, rdb, cfg.ScyllaKeyspace)
	}

	repo := postgres.NewPostgresRepository(dbpool)
	authUseCase := usecase.NewAuthUseCase(repo, cfg.JWTSecret, cfg.JWTTTL)
	userUseCase := usecase.NewUserUseCase(repo)

	hubCtx, hubCancel := context.WithCancel(ctx)
	defer hubCancel()

	hub := ws.NewHub(msgRepo)
	go hub.Run(hubCtx)

	h := server.Default(
		server.WithHostPorts(cfg.Addr),
		server.WithHandleMethodNotAllowed(true),
	)

	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		hub.Stop()
		hubCancel()
	})

	upgrader := ws.NewUpgrader(cfg.AllowedOrigins)
	wsHandler := ws.ServeWs(hub, upgrader)

	cookieCfg := http.CookieConfig{
		Name:   "token",
		MaxAge: int(cfg.JWTTTL.Seconds()),
		Secure: false,
		Domain: "",
	}

	http.SetupRouter(h, http.Deps{
		Auth:      authUseCase,
		Users:     userUseCase,
		MsgRepo:   msgRepo,
		WsHandler: wsHandler,
		JWTSecret: cfg.JWTSecret,
		Cookie:    cookieCfg,
	})

	h.Spin()
	return nil
}
