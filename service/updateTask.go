package service

import (
	"2do.com/models"
	"2do.com/repository"
	"context"
	"fmt"
)

func UpdateTask(ctx context.Context, taskRepository repository.TaskRepository, task models.Task) (*models.Task, error) {
	fmt.Println("Task updated")
	updatedTask, err := taskRepository.UpdateTask(ctx, task.ID, task)
	if err != nil {
		return nil, err
	}
	return updatedTask, nil
}
