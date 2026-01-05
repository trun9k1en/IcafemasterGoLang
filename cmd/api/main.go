package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"icafe-registration/internal/config"
	httpDelivery "icafe-registration/internal/delivery/http"
	"icafe-registration/internal/domain"
	"icafe-registration/internal/repository/mongodb"
	"icafe-registration/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongoDB, err := config.NewMongoDB(&cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ensure upload directories exist
	if err := os.MkdirAll(cfg.Upload.Path+"/files", 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	if err := os.MkdirAll(cfg.Upload.Path+"/videos", 0755); err != nil {
		log.Fatalf("Failed to create videos directory: %v", err)
	}

	// Initialize repositories
	registrationRepo := mongodb.NewRegistrationRepository(mongoDB.Database)
	fileRepo := mongodb.NewFileRepository(mongoDB.Database)
	userRepo := mongodb.NewUserRepository(mongoDB.Database)

	// Initialize usecases
	contextTimeout := 10 * time.Second
	registrationUsecase := usecase.NewRegistrationUsecase(registrationRepo, contextTimeout)
	fileUsecase := usecase.NewFileUsecase(fileRepo, &cfg.Upload, contextTimeout)
	authUsecase := usecase.NewAuthUsecase(userRepo, &cfg.JWT, contextTimeout)
	userUsecase := usecase.NewUserUsecase(userRepo, contextTimeout)

	// Create default admin user if not exists
	createDefaultAdmin(userRepo, userUsecase)

	// Initialize HTTP router
	router := httpDelivery.NewRouter(
		registrationUsecase,
		fileUsecase,
		authUsecase,
		userUsecase,
		cfg,
	)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := router.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Close MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoDB.Close(ctx); err != nil {
		log.Printf("Error closing MongoDB connection: %v", err)
	}

	log.Println("Server exited gracefully")
}

// createDefaultAdmin creates a default admin user if not exists
func createDefaultAdmin(userRepo domain.UserRepository, userUsecase domain.UserUsecase) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if admin exists
	_, err := userRepo.GetByUsername(ctx, "admin")
	if err == nil {
		log.Println("Default admin user already exists")
		return
	}

	// Create default admin
	req := &domain.CreateUserRequest{
		Username: "admin",
		Email:    "admin@icafe.local",
		Password: "admin123", // Change this in production!
		FullName: "Administrator",
		Role:     domain.RoleAdmin,
	}

	_, err = userUsecase.Create(ctx, req)
	if err != nil {
		log.Printf("Failed to create default admin: %v", err)
		return
	}

	log.Println("Default admin user created successfully")
	log.Println("Username: admin")
	log.Println("Password: admin123 (Please change this in production!)")
}
