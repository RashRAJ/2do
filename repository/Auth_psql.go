package repository

import (
	"2do.com/models"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type UserRepository struct {
	*BaseRepository
}

func NewAuthRepository(db *pgxpool.Pool, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		BaseRepository: &BaseRepository{
			db:     db,
			logger: logger,
		},
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	start := time.Now()
	r.logger.Info("creating user", zap.String("username", user.Username))

	err := r.db.QueryRow(ctx, `INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`,
		user.Username, user.Email, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, pgxError) {
			if pgxError.Code == "23505" {
				r.logger.Warn("duplicate user creation attempted",
					zap.String("username", user.Username),
					zap.String("postgres_code", pgxError.Code),
					zap.Duration("duration", time.Since(start)),
				)
				return nil, ErrDuplicate
			}
		}
		r.logger.Error("failed to insert user",
			zap.Error(err),
			zap.String("username", user.Username),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}
	r.logger.Info("user created successfully",
		zap.Int("user_id", user.ID),
		zap.String("username", user.Username),
		zap.Duration("duration", time.Since(start)),
	)
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	start := time.Now()
	r.logger.Info("retrieving user by username", zap.String("username", username))

	user := &models.User{}
	err := r.db.QueryRow(ctx, `SELECT id, username, email, password_hash, created_at FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		r.logger.Error("failed to query user by username",
			zap.Error(err),
			zap.String("username", username),
			zap.Duration("duration", time.Since(start)),
		)
		return *user, err
	}
	r.logger.Info("retrieved user by username successfully")
	return *user, nil
}
