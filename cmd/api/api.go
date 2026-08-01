package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	_ "megan-messenger/docs"
	"megan-messenger/internal/auth"
	"megan-messenger/internal/config"
	"megan-messenger/internal/conversation"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/message"
	"megan-messenger/internal/user"
	"megan-messenger/internal/ws"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.config.CORS.AllowedOrigins,
		AllowedMethods:   app.config.CORS.AllowedMethods,
		AllowedHeaders:   app.config.CORS.AllowedHeaders,
		ExposedHeaders:   app.config.CORS.ExposedHeaders,
		AllowCredentials: app.config.CORS.AllowCredentials,
		MaxAge:           app.config.CORS.MaxAge,
	}))

	r.Get("/docs/*", httpSwagger.Handler())

	authMw := auth.Middleware(app.authenticator)
	wsAuthMw := auth.WebSocketMiddleware(app.authenticator)

	authHandler := auth.NewHandler(app.validator, app.authService)
	userHandler := user.NewHandler(app.validator, app.userService)
	conversationHandler := conversation.NewHandler(app.validator, app.conversationService, app.wsService)
	messageHandler := message.NewHandler(app.validator, app.messageService)

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/phone/start", authHandler.StartPhoneAuth)
		r.Post("/phone/verify", authHandler.VerifyPhone)
		r.Post("/password/verify", authHandler.VerifyPasswordChallenge)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)

		r.With(authMw).Post("/onboarding/username", authHandler.CompleteUsername)
		r.With(authMw, auth.RequireOnboarded).Post("/password", authHandler.SetPassword)
		r.With(authMw, auth.RequireOnboarded).Delete("/password", authHandler.RemovePassword)
	})

	r.With(authMw).Group(func(r chi.Router) {
		r.Route("/users/me", func(r chi.Router) {
			r.Get("/", userHandler.CurrentUser)
			r.With(auth.RequireOnboarded).Post("/avatar", userHandler.UploadAvatar)
		})

		r.With(auth.RequireOnboarded).Route("/conversations", func(r chi.Router) {
			r.Get("/", conversationHandler.List)
			r.Post("/group", conversationHandler.CreateGroup)
			r.Post("/dm", conversationHandler.CreateOrGetDM)
			r.Post("/join/{slug:[a-z0-9]{11}}", conversationHandler.JoinBySlug)
			r.Route("/{conversationID}", func(r chi.Router) {
				r.Get("/messages", messageHandler.ListMessages)
			})
		})
	})

	r.With(wsAuthMw, auth.RequireOnboarded).
		Get("/conversations/{conversationID}/ws", conversationHandler.Websocket)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Addr),
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	slog.Info("server has started", "addr", app.config.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	slog.Info("completing background tasks", "addr", srv.Addr)

	app.db.Close()
	if err := app.rdb.Close(); err != nil {
		slog.Error("failed to close redis", "error", err)
	}

	slog.Info("server stopped", "addr", srv.Addr)

	return nil
}

type application struct {
	config              *config.Config
	db                  *pgxpool.Pool
	rdb                 *redis.Client
	authenticator       *auth.Authenticator
	authService         *auth.Service
	userService         *user.Service
	conversationService *conversation.Service
	wsService           *ws.Service
	messageService      *message.Service
	validator           *httputil.Validator
}
