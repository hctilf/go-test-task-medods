package routes

import (
	app "github.com/hctilf/go-test-task-medods/internal/app"

	"github.com/hctilf/go-test-task-medods/internal/controller/http/handlers/auth"
	mw "github.com/hctilf/go-test-task-medods/internal/controller/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func SetRoutes(app *app.Application, server *fiber.App) error {
	server.Use(
		requestid.New(requestid.ConfigDefault),
		mw.Logger(app),
		logger.New(
			logger.Config{
				Format: "[${time}] ${referer} ${ip}:${port} ${status} ${latency} - ${method} ${path}\n",
			},
		),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		recover.New(),
	)

	authHandler := auth.NewAuthHandler(app)

	api := server.Group("/api")
	api.Route("/auth", func(api fiber.Router) {
		api.Get("/tokens", authHandler.GetAccessNRefreshToken)
		api.Post("/refresh", authHandler.PostRefreshToken)
		api.Post("/test", authHandler.TestTokenCreate)
	})

	return nil
}

func SetSession(sessionOnly bool) *session.Store {
	return session.New(session.Config{
		KeyLookup:         "cookie:__Host-session",
		CookieSecure:      true,
		CookieHTTPOnly:    true,
		CookieSameSite:    "Strict",
		CookieSessionOnly: sessionOnly,
	})
}
