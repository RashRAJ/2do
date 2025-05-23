package service

import (
	"2do.com/repository"
	"context"
	"fmt"
)

func Migrate(ctx context.Context, taskRepository repository.TaskRepository) error {
	fmt.Println("Migration happening now")
	return taskRepository.Migrate(ctx)
}
