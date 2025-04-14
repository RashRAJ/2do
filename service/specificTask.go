package service

import (
	"errors"
	"log/slog"
	"time"

	"../models"
	"../repository"
	"go.uber.org/zap"
)

type TaskService struct {
	repo   *repository.TaskRepository
	logger *zap.Logger
}

func NewTaskService(repo *repository.TaskRepository, logger *zap.Logger) *TaskService {
	return &TaskService{repo: repo, logger: logger}
}

func SpecificTask(id int) (*models.Task, error) {
	if id <= 0 {
		s.logger.Error("Invalid task ID", zap.Int("id", id))
		return nil, errors.New("invalid task ID")
	}
	task, err := s.repo.SpecificTask(id)
	if err != nil {
		s.logger.Error("Failed to fetch specific task", zap.int("id", id), zap.Error(err))
		return nil, err
	}
	s.Logger.info("Fetch task sucessfully", zap.int("id", id))
	return task, nil
}

//func specificTask(w http.ResponseWriter, r *http.Request) {
//	initlogger()
//	params := mux.Vars(r)
//	id, err := strconv.Atoi(params["id"])
//	if err != nil {
//		http.Error(w, "Invalid ID format", http.StatusBadRequest)
//		logger.Error("Invalid ID formart", zap.Error(err))
//		return
//	}
//
//	var task Task
//	var createdBytes []byte // Temporary variable to hold the raw datetime value
//
//	// Query the database and scan the result
//	err = db.QueryRow("SELECT id, title, content, created, status FROM tasks WHERE id = ?", id).
//		Scan(&task.ID, &task.Title, &task.Content, &createdBytes, &task.Status)
//	if err == sql.ErrNoRows {
//		http.Error(w, "Task not found", http.StatusNotFound)
//		logger.Error("Task not found", zap.Error(err))
//		return
//	}
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Error("Failed to query db", zap.Error(err))
//		return
//	}
//
//	// Parse the raw datetime value into a time.Time object
//	created, err := time.Parse("2006-01-02 15:04:05", string(createdBytes))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Error("Failed to parse datetime", zap.Error(err))
//		return
//	}
//	task.Created = created
//
//	// Encode the task as JSON and send it in the response
//	json.NewEncoder(w).Encode(task)
//}
