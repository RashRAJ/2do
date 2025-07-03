package repository

import (
	"2do.com/models"
	"context"
	"errors"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Task) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	DeleteTask(ctx context.Context, id int) error
	UpdateTask(ctx context.Context, id int, updated models.Task) (*models.Task, error)
	GetbyID(ctx context.Context, id int) (models.Task, error)
}

// interface got implemented here to be used by the service layer
