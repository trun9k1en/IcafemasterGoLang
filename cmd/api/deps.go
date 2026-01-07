package main

import (
	"icafe-registration/internal/config"
	"icafe-registration/internal/repository/mongodb"
	"icafe-registration/internal/usecase"
	"time" // Cần import time để sử dụng Duration
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
		Customer:     mongodb.NewCustomerRepository(a.Database.MongoDB.Database),
	}
}

// initUsecases initializes all usecases
func (a *App) initUsecases() {
	// 1. Khai báo contextTimeout (Lấy từ config hoặc set mặc định)
	// Bạn có thể dùng: contextTimeout := time.Duration(a.Config.App.ContextTimeout) * time.Second
	contextTimeout := 10 * time.Second

	a.Usecases = &UsecaseDeps{
		// 2. CẬP NHẬT: Truyền thêm a.Repos.Customer vào NewRegistrationUsecase
		Registration: usecase.NewRegistrationUsecase(
			a.Repos.Registration,
			a.Repos.Customer, // Thêm tham số này để lưu data vào bảng customers
			contextTimeout,
		),

		File:     usecase.NewFileUsecase(a.Repos.File, &a.Config.Upload, contextTimeout),
		Auth:     usecase.NewAuthUsecase(a.Repos.User, &a.Config.JWT, contextTimeout),
		User:     usecase.NewUserUsecase(a.Repos.User, contextTimeout),
		Customer: usecase.NewCustomerUsecase(a.Repos.Customer, contextTimeout),
	}
}
