package repository

import (
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
	Migrate(ctx context.Context) error
	CreateTask(ctx context.Context, task Task) (*Task, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	DeleteTask(ctx context.Context, id int) error
	UpdateTask(ctx context.Context, id int, updated Task) (*Task, error)
	GetbyID(ctx context.Context, id int) (Task, error)
}
