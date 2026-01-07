package main

import (
	"context"
	"log"
	"os"
	"time"

	"icafe-registration/internal/config"
	httpDelivery "icafe-registration/internal/delivery/http"
)

const contextTimeout = 10 * time.Second

// NewApp creates and initializes a new application
func NewApp() (*App, error) {
	app := &App{}

	app.Config = config.LoadConfig()

	if err := app.initDatabase(); err != nil {
		return nil, err
	}

	if err := app.initDirectories(); err != nil {
		return nil, err
	}

	app.initRepositories()
	app.initUsecases()
	app.createDefaultUsers()
	app.initRouter()

	return app, nil
}

// initDirectories creates necessary upload directories
func (a *App) initDirectories() error {
	dirs := []string{
		a.Config.Upload.Path + "/files",
		a.Config.Upload.Path + "/videos",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// initRouter initializes HTTP router
func (a *App) initRouter() {
	a.Router = httpDelivery.NewRouter(
		a.Usecases.Registration,
		a.Usecases.File,
		a.Usecases.Auth,
		a.Usecases.User,
		a.Usecases.Customer,
		a.Config,
	)
}

// Run starts the application
func (a *App) Run() error {
	log.Printf("Server starting on %s:%s", a.Config.Server.Host, a.Config.Server.Port)
	return a.Router.Run()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	if err := a.Database.MongoDB.Close(ctx); err != nil {
		log.Printf("Error closing MongoDB connection: %v", err)
		return err
	}

	log.Println("Server exited gracefully")
	return nil
}
