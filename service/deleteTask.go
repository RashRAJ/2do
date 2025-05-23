package service

import (
	"2do.com/repository"
	"context"
	"fmt"
)

func DeleteTask(ctx context.Context, taskRepository repository.TaskRepository, id int) error {
	fmt.Println("Task deleted")
	return taskRepository.DeleteTask(ctx, id)
}
