package main

import (
	"context"
	"log/slog"
	"megan-messenger/internal/auth"
	"megan-messenger/internal/config"
	"megan-messenger/internal/conversation"
	"megan-messenger/internal/db/postgres"
	db "megan-messenger/internal/db/postgres/sqlc"
	redis2 "megan-messenger/internal/db/redis"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/message"
	"megan-messenger/internal/notification"
	"megan-messenger/internal/user"
	"megan-messenger/internal/ws"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// @title						Megan Messenger API
// @version					1.0.0
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed load config", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})

	queries := db.New(pool)

	accessCfg := cfg.Auth.AccessToken
	authenticator := auth.NewJWTAuthenticator(accessCfg.Secret, accessCfg.Issuer, accessCfg.TTL)
	refreshCfg := cfg.Auth.RefreshToken
	refreshRepo := redis2.NewRefreshTokenRepository(rdb, refreshCfg.KeyPrefix, refreshCfg.UserKeyPrefix, refreshCfg.TTL)
	challengeRepo := redis2.NewPasswordChallengeRepository(rdb, cfg.Auth.Challenge.KeyPrefix, cfg.Auth.Challenge.TTL)
	userRepo := postgres.NewUserRepository(queries)
	conversationRepo := postgres.NewConversationRepository(queries)
	messageRepo := postgres.NewMessageRepository(queries)

	authService := auth.NewService(
		authenticator,
		refreshRepo,
		challengeRepo,
		userRepo,
		notification.NewLogSmsSender(logger),
		cfg.Auth,
	)
	userService := user.NewService(userRepo, cfg.FileStore.AvatarsPath())
	conversationService := conversation.NewService(conversationRepo)
	wsService := ws.NewService(rdb, userRepo, messageRepo, cfg.CORS.AllowedOrigins)
	messageService := message.NewService(conversationRepo, messageRepo)
	validator := httputil.NewValidator()

	api := application{
		config:              cfg,
		db:                  pool,
		rdb:                 rdb,
		authenticator:       authenticator,
		authService:         authService,
		userService:         userService,
		conversationService: conversationService,
		wsService:           wsService,
		messageService:      messageService,
		validator:           validator,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
