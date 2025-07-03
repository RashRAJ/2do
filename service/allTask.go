package service

import (
	"2do.com/middleware"
	"2do.com/models"
	"2do.com/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
)

func (s *TaskService) GetAllTasks(ctx context.Context, taskRepository repository.TaskRepository) ([]models.Task, error) {
	fmt.Println("GetAllTasks happening now")
	logger := middleware.ZapLogger
	tasks, err := taskRepository.GetAllTasks(ctx)
	if err != nil {
		logger.Error("Error getting all tasks", zap.Error(err))
		return nil, err
	}
	logger.Info("Tasks retrieved successfully")
	fmt.Printf("Service retrieved %d tasks\n", len(tasks))
	return tasks, nil
}
