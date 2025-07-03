package service

import (
	"2do.com/middleware"
	"2do.com/models"
	"2do.com/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
)

func (s *TaskService) GetTaskByID(ctx context.Context, taskRepository repository.TaskRepository, id int) (*models.Task, error) {
	fmt.Println("Single Task")
	logger := middleware.ZapLogger
	task, err := taskRepository.GetbyID(ctx, id)
	if err != nil {
		logger.Error("Error getting task by ID", zap.Int("task_id", id))
		return nil, err
	}
	return &task, nil
}
