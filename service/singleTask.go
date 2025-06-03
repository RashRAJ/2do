package service

import (
	"2do.com/models"
	"2do.com/repository"
	"context"
	"fmt"
)

func GetTaskByID(ctx context.Context, taskRepository repository.TaskRepository, id int) (*models.Task, error) {
	fmt.Println("Single Task")
	task, err := taskRepository.GetbyID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
