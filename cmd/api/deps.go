package main

import (
	"icafe-registration/internal/config"
	"icafe-registration/internal/repository/mongodb"
	"icafe-registration/internal/usecase"
)

// initDatabase initializes database connections
func (a *App) initDatabase() error {
	mongoDB, err := config.NewMongoDB(&a.Config.MongoDB)
	if err != nil {
		return err
	}

	a.Database = &DatabaseDeps{
		MongoDB: mongoDB,
	}

	return nil
}

// initRepositories initializes all repositories
func (a *App) initRepositories() {
	a.Repos = &RepositoryDeps{
		Registration: mongodb.NewRegistrationRepository(a.Database.MongoDB.Database),
		File:         mongodb.NewFileRepository(a.Database.MongoDB.Database),
		User:         mongodb.NewUserRepository(a.Database.MongoDB.Database),
	}
}

// initUsecases initializes all usecases
func (a *App) initUsecases() {
	a.Usecases = &UsecaseDeps{
		Registration: usecase.NewRegistrationUsecase(a.Repos.Registration, contextTimeout),
		File:         usecase.NewFileUsecase(a.Repos.File, &a.Config.Upload, contextTimeout),
		Auth:         usecase.NewAuthUsecase(a.Repos.User, &a.Config.JWT, contextTimeout),
		User:         usecase.NewUserUsecase(a.Repos.User, contextTimeout),
	}
}
