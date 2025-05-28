package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"
	"path/filepath"
	"strings"
	"web/api/v1"

	"web/config"
	_ "web/docs"
	"web/middleware"
	"web/repos"
	"web/services"
)

// @title Course API
// @version 1.0
// @description This is a course management API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("Starting application")

	// Load application configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Info("Configuration loaded successfully")

	// Run SQL migrations if needed
	if err := runMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Info("Migrations applied successfully")

	// Initialize repositories
	courseRepo := repos.NewCourseRepository(appConfig.GormDB)
	chapterRepo := repos.NewChapterRepository(appConfig.GormDB)
	lessonRepo := repos.NewLessonRepository(appConfig.GormDB)
	userRepo := repos.NewUserRepository(appConfig.GormDB)

	// Initialize services
	courseService := services.NewCourseService(courseRepo)
	chapterService := services.NewChapterService(chapterRepo, courseRepo)
	lessonService := services.NewLessonService(lessonRepo, chapterRepo, courseRepo)
	authService := services.NewAuthService(appConfig, userRepo)
	userService := services.NewUserService(userRepo)

	// Initialize router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.ResponseMiddleware())

	// Apply auth middleware to all routes except swagger, root, and auth endpoints
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") ||
			c.Request.URL.Path == "/" ||
			strings.HasPrefix(c.Request.URL.Path, "/api/v1/auth/") {
			c.Next()
			return
		}

		// Apply authentication middleware
		middleware.AuthMiddleware(authService)(c)
	})

	// Register api
	courseHandler := v1.NewCourseHandler(appConfig, courseService, chapterService, authService)
	chapterHandler := v1.NewChapterHandler(appConfig, chapterService, authService)
	lessonHandler := v1.NewLessonHandler(appConfig, lessonService, authService)
	userHandler := v1.NewUserHandler(appConfig, userService, authService)

	// Register routes
	courseHandler.RegisterRoutes(router)
	chapterHandler.RegisterRoutes(router)
	lessonHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)

	// Default route
	router.GET("/", func(c *gin.Context) {
		middleware.RespondWithSuccess(c, nil, "Hello World!")
	})

	// Swagger documentation endpoint
	url := httpSwagger.URL("/swagger/v1/doc.json") // The URL pointing to API definition
	router.GET("/swagger/v1/*any", gin.WrapH(httpSwagger.Handler(url)))

	log.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	migrationsDir := filepath.Join(dir, "migrations")

	// Load config to get DB connection
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get the SQL DB from AppConfig
	db := appConfig.DB

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
