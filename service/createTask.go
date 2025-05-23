package service

import (
	"2do.com/repository"
	"context"
)

func CreateTask(ctx context.Context, taskRepository repository.TaskRepository, task repository.Task) (*repository.Task, error) {
	createdTask, err := taskRepository.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return createdTask, nil
}
