package repository

import (
	"2do.com/models"
	"context"
	"errors"
)

var (
	ErrDuplicate    = errors.New("user already exists")
	ErrNotExist     = errors.New("user does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, email string) error
}
