package app

import (
	"sync"

	config "github.com/hctilf/go-test-task-medods/config"

	"go.uber.org/zap"

	uc "github.com/hctilf/go-test-task-medods/internal/usecase"

	"github.com/gofiber/fiber/v2/middleware/session"
)

type Application struct {
	Config  *config.Config
	Logger  *zap.SugaredLogger
	Tokens  *uc.TokensUsecase
	Session *session.Store
	WG      sync.WaitGroup
}

func (app *Application) Background(fn func()) {
	app.WG.Add(1)
	go func() {
		defer app.WG.Done()

		defer func() {
			if err := recover(); err != nil {
				app.Logger.Error("background error:", err)
			}
		}()

		fn()
	}()
}
