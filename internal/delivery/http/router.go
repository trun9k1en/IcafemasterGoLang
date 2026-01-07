package http

import (
	"icafe-registration/internal/config"
	"icafe-registration/internal/domain"

	"github.com/gin-gonic/gin"
)

// Router holds all dependencies for HTTP router
type Router struct {
	Engine              *gin.Engine
	RegistrationUsecase domain.RegistrationUsecase
	FileUsecase         domain.FileUsecase
	AuthUsecase         domain.AuthUsecase
	UserUsecase         domain.UserUsecase
	CustomerUsecase     domain.CustomerUsecase
	Config              *config.Config
}

// NewRouter creates a new HTTP router
func NewRouter(
	registrationUsecase domain.RegistrationUsecase,
	fileUsecase domain.FileUsecase,
	authUsecase domain.AuthUsecase,
	userUsecase domain.UserUsecase,
	customerUsecase domain.CustomerUsecase,
	cfg *config.Config,
) *Router {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// Apply middlewares
	engine.Use(LoggerMiddleware())
	engine.Use(RecoveryMiddleware())
	engine.Use(CORSMiddleware())

	// Set max multipart memory for file uploads
	engine.MaxMultipartMemory = cfg.Upload.MaxFileSize

	router := &Router{
		Engine:              engine,
		RegistrationUsecase: registrationUsecase,
		FileUsecase:         fileUsecase,
		AuthUsecase:         authUsecase,
		UserUsecase:         userUsecase,
		CustomerUsecase:     customerUsecase,
		Config:              cfg,
	}

	router.setupRoutes()

	return router
}

// setupRoutes sets up all routes
func (r *Router) setupRoutes() {
	// Health check
	r.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := r.Engine.Group("/api/v1")
	{
		// Public routes - Auth
		NewAuthHandler(v1, r.AuthUsecase)

		// Public routes - Registration (anyone can register)
		NewRegistrationHandler(v1, r.RegistrationUsecase)

		// Public file serving routes
		NewFileHandler(v1, r.Engine, r.FileUsecase, &r.Config.Upload)

		// Protected routes - require authentication
		protected := v1.Group("")
		protected.Use(JWTAuthMiddleware(r.AuthUsecase))
		{
			// User management routes (admin only)
			adminOnly := protected.Group("")
			adminOnly.Use(RequireRole(domain.RoleAdmin))
			{
				NewUserHandler(adminOnly, r.UserUsecase)
			}

			// Customer routes (admin can CRUD, sale can only read)
			NewCustomerHandler(protected, r.CustomerUsecase)
		}
	}
}

// Run starts the HTTP server
func (r *Router) Run() error {
	addr := r.Config.Server.Host + ":" + r.Config.Server.Port
	return r.Engine.Run(addr)
}
