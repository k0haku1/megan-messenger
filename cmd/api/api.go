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
	"megan-messenger/internal/project"
	"megan-messenger/internal/ratelimit"
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
	userHandler := user.NewHandler(app.validator, app.userService, app.rateLimiter)
	conversationHandler := conversation.NewHandler(app.validator, app.conversationService, app.wsService)
	messageHandler := message.NewHandler(app.validator, app.messageService, app.wsService)
	projectHandler := project.NewHandler(app.validator, app.projectService)

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
		r.Get("/users/me", userHandler.CurrentUser)

		r.With(auth.RequireOnboarded).Group(func(r chi.Router) {
			r.Get("/users/search", userHandler.SearchUsers)
			r.Get("/users/by-username/{username}", userHandler.GetByUsername)
			r.Patch("/users/me/username", userHandler.UpdateUsername)
			r.Patch("/users/me/privacy", userHandler.UpdatePrivacy)
			r.Post("/users/me/avatar", userHandler.UploadAvatar)

			r.Route("/conversations", func(r chi.Router) {
				r.Get("/", conversationHandler.List)
				r.Post("/group", conversationHandler.CreateGroup)
				r.Post("/dm/messages", conversationHandler.SendDMMessage)
				r.Post("/join/{slug:[a-z0-9]{11}}", conversationHandler.JoinBySlug)
				r.Route("/{conversationID}", func(r chi.Router) {
					r.Post("/leave", conversationHandler.LeaveGroup)
					r.Get("/projects", projectHandler.ListConversationProjects)
					r.Get("/messages", messageHandler.ListMessages)
					r.Post("/messages", messageHandler.SendMessage)
				})
			})

			r.Route("/projects", func(r chi.Router) {
				r.Get("/", projectHandler.List)
				r.Post("/", projectHandler.Create)
				r.Route("/{projectID}", func(r chi.Router) {
					r.Get("/", projectHandler.Get)
					r.Get("/members", projectHandler.ListMembers)
					r.Post("/members", projectHandler.AddMember)
					r.Get("/docs", projectHandler.ListDocs)
					r.Post("/docs", projectHandler.CreateDoc)
					r.Route("/docs/{docID}", func(r chi.Router) {
						r.Get("/", projectHandler.GetDoc)
						r.Patch("/", projectHandler.UpdateDoc)
						r.Delete("/", projectHandler.DeleteDoc)
					})
					r.Get("/decisions", projectHandler.ListDecisions)
					r.Post("/decisions", projectHandler.CreateDecision)
					r.Patch("/decisions/{decisionID}", projectHandler.UpdateDecisionStatus)
					r.Get("/conversations", projectHandler.ListConversations)
					r.Post("/conversations", projectHandler.LinkConversation)
					r.Delete("/conversations/{conversationID}", projectHandler.UnlinkConversation)
				})
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
	projectService      *project.Service
	validator           *httputil.Validator
	rateLimiter         *ratelimit.Limiter
}
