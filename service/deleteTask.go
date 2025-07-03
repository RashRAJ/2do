package service

import (
	"2do.com/middleware"
	"2do.com/repository"
	"context"
	"go.uber.org/zap"
)

func (j *TaskService) DeleteTask(ctx context.Context, taskRepository repository.TaskRepository, id int) error {
	logger := middleware.ZapLogger
	logger.Info("deleting task", zap.Int("task_id", id))
	return taskRepository.DeleteTask(ctx, id)
}
