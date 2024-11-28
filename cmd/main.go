package main

import (
	"encoding/gob"
	"os"
	"os/signal"
	"sync"
	"syscall"

	app "github.com/hctilf/go-test-task-medods/internal/app"

	"github.com/google/uuid"

	config "github.com/hctilf/go-test-task-medods/config"
	logger "github.com/hctilf/go-test-task-medods/pkg/logger"

	postgres "github.com/hctilf/go-test-task-medods/pkg/postgres"

	entity "github.com/hctilf/go-test-task-medods/internal/entity"
	uc "github.com/hctilf/go-test-task-medods/internal/usecase"
	tokens "github.com/hctilf/go-test-task-medods/internal/usecase/repo/tokens"

	routes "github.com/hctilf/go-test-task-medods/internal/controller/http/routes"
	server "github.com/hctilf/go-test-task-medods/internal/controller/http/server"
)

var (
	appInstance *app.Application
	once        sync.Once
)

func init() {
	gob.Register(uuid.UUID{})
}

func main() {
	_ = server.GenerateSelfSignedCert("cert.pem", "key.pem")
	app := GetApp()
	server := server.NewServer(app)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-interrupt:
		app.Logger.Errorf("signal.Notify: %v", s)
	case err := <-server.Notify():
		app.Logger.Errorf("server.Notify: %v", err)
	}

	app.WG.Wait()
	err := server.Stop()
	if err != nil {
		app.Logger.Errorf("server.Stop: %v", err)
	}
}

func GetApp() *app.Application {
	once.Do(func() {
		appInstance = newApplication()
	})

	return appInstance
}

func newApplication() *app.Application {
	config := config.GetConfig()
	pgStorage := postgres.GetDB(config)
	_ = postgres.DoMigrate(pgStorage, &entity.RefreshToken{})
	logger := logger.NewLogger(config.Env.LogLevel)

	return &app.Application{
		Config:  config,
		Logger:  logger,
		Tokens:  uc.NewTokensUsecase(tokens.NewTokensRepository(pgStorage)),
		Session: routes.SetSession(false),
	}
}
