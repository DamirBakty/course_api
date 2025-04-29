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

	return &AppConfig{
		DB:     sqlDB,
		GormDB: gormDB,
		DbUrl:  dbUrl,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
