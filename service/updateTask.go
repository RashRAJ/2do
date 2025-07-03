package service

import (
	"2do.com/middleware"
	"2do.com/models"
	"2do.com/repository"
	"context"
	"go.uber.org/zap"
)

func (s *TaskService) UpdateTask(ctx context.Context, taskRepository repository.TaskRepository, task models.Task) (*models.Task, error) {
	logger := middleware.ZapLogger
	updatedTask, err := taskRepository.UpdateTask(ctx, task.ID, task)
	if err != nil {
		logger.Error("Error updating task", zap.Error(err))
		return nil, err
	}
	return updatedTask, nil
}
