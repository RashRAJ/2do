package db

import (
	"2do.com/config"
	"2do.com/middleware"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	logger := middleware.ZapLogger
	dbConfig := config.LoadDBConfig()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)

	m, err := migrate.New(
		"file://migrations", dbURL)

	if err != nil {
		logger.Error("failed to initialize migration", zap.Error(err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Migration failed: %v", zap.Error(err))
	}
	logger.Info("Migration completed")
}
