package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/config"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/domain"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/handler"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/middleware"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/repository/postgres"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/service"
	customJWT "github.com/youngermaster/my-distributed-bookstore/services/users-service/pkg/jwt"
	postgresql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Initialize logger
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if os.Getenv("ENV") == "development" {
		log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting Users Service...")

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := connectDB(cfg.Database, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Seed default roles
	if err := service.SeedRoles(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to seed roles")
	}

	// Initialize JWT manager
	jwtManager := customJWT.NewJWTManager(cfg.JWT.Secret, cfg.JWT.GetTokenDuration())

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	wishlistRepo := postgres.NewWishlistRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	wishlistService := service.NewWishlistService(wishlistRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	wishlistHandler := handler.NewWishlistHandler(wishlistService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger(log))

	// Health check endpoints
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		// Check database connection
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  "database connection error",
			})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  "database ping failed",
			})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes (require authentication)
	auth.Post("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)

	// User profile routes (protected)
	users := api.Group("/users", middleware.AuthMiddleware(authService))
	users.Get("/me", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		userEmail := c.Locals("userEmail")
		return c.JSON(fiber.Map{
			"id":    userID,
			"email": userEmail,
		})
	})

	// Wishlist routes (protected)
	users.Get("/me/wishlist", wishlistHandler.GetWishlist)
	users.Post("/me/wishlist", wishlistHandler.AddToWishlist)
	users.Delete("/me/wishlist/:book_id", wishlistHandler.RemoveFromWishlist)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Info().Str("addr", addr).Msg("Users Service listening")
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	log.Info().Msg("Server stopped")
}

func connectDB(cfg config.DatabaseConfig, log zerolog.Logger) (*gorm.DB, error) {
	log.Info().Msg("Connecting to database...")

	db, err := gorm.Open(postgresql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().Msg("Database connected successfully")
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.UserRole{},
		&domain.Address{},
		&domain.Session{},
		&domain.WishlistItem{},
	)
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
