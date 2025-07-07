package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// BaseRepository contains shared dependencies
type BaseRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}
