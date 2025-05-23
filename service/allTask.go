package service

import (
	"2do.com/repository"
	"context"
	"fmt"
)

func GetAllTasks(ctx context.Context, taskRepository repository.TaskRepository) ([]repository.Task, error) {
	fmt.Println("GetAllTasks happening now")
	tasks, err := taskRepository.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
