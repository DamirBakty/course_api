package config

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB     *sql.DB
	GormDB *gorm.DB
	DbUrl  string

	// Keycloak configuration
	KeycloakURL           string
	KeycloakRealm         string
	KeycloakClientID      string
	KeycloakClientSecret  string
	KeycloakAdminUsername string
	KeycloakAdminPassword string

	// MinIO configuration
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbUrl := getEnv("DB_URL", "postgres://postgres:password@localhost:5432/database?sslmode=disable")

	// Initialize standard SQL DB for migrations
	sqlDB, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Initialize GORM DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Load Keycloak configuration
	keycloakURL := getEnv("KEYCLOAK_URL", "http://localhost:8081")
	keycloakRealm := getEnv("KEYCLOAK_REALM", "master")
	keycloakClientID := getEnv("KEYCLOAK_CLIENT_ID", "course-api")
	keycloakClientSecret := getEnv("KEYCLOAK_CLIENT_SECRET", "")
	keycloakAdminUsername := getEnv("KC_ADMIN", "admin")
	keycloakAdminPassword := getEnv("KC_ADMIN_PASSWORD", "admin")

	// Load MinIO configuration
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getEnv("MINIO_BUCKET", "attachments")
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	return &AppConfig{
		DB:                    sqlDB,
		GormDB:                gormDB,
		DbUrl:                 dbUrl,
		KeycloakURL:           keycloakURL,
		KeycloakRealm:         keycloakRealm,
		KeycloakClientID:      keycloakClientID,
		KeycloakClientSecret:  keycloakClientSecret,
		KeycloakAdminUsername: keycloakAdminUsername,
		KeycloakAdminPassword: keycloakAdminPassword,
		MinioEndpoint:         minioEndpoint,
		MinioAccessKey:        minioAccessKey,
		MinioSecretKey:        minioSecretKey,
		MinioBucket:           minioBucket,
		MinioUseSSL:           minioUseSSL,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
