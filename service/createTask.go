package service

import (
	"2do.com/middleware"
	"2do.com/models"
	"2do.com/repository"
	"context"
	"go.uber.org/zap"
)

func (s *TaskService) CreateTask(ctx context.Context, taskRepository repository.TaskRepository, task models.Task) (*models.Task, error) {
	createdTask, err := taskRepository.CreateTask(ctx, task)
	logger := middleware.ZapLogger
	if err != nil {
		logger.Error("Error creating task", zap.Error(err))
		return nil, err
	}
	return createdTask, nil
}
