package service

import (
	"2do.com/repository"
	"go.uber.org/zap"
)

type TaskService struct {
	repo   *repository.TaskRepository
	logger *zap.Logger
}

func NewTaskService(repo *repository.TaskRepository, logger *zap.Logger) *TaskService {
	return &TaskService{repo: repo, logger: logger}
}
