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
		fmt.Printf("Error getting all tasks: %v\n", err)
		return nil, err
	}
	fmt.Printf("Service retrieved %d tasks\n", len(tasks))
	return tasks, nil
}
