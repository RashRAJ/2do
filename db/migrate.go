package db

import (
	"2do.com/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigrations() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		logger = zap.NewNop() // Safe fallback
	}
	defer logger.Sync()

	dbConfig := config.LoadDBConfig()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)

	m, err := migrate.New(
		"file://db/migrations", dbURL)

	if err != nil {
		logger.Error("failed to initialize migration", zap.Error(err))
		return
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Migration failed ", zap.Error(err))
		return
	}
	logger.Info("Migration completed")
}
