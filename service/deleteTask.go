package service

import (
	"errors"
	"time"

	"../models"
	"../repository/"
	"go.uber.org/zap"
)

type TaskService struct {
	repo   *repository.TaskRepository
	logger *zap.Logger
}

func (s *TaskService) DeleteTask(id int) error {
	if id <= 0 {
		return errors.New("invalid task ID")
	}

	rowsAffected, err := s.repo.DeleteTask(id)
	if err != nil {
		s.logger.Error("Repository error deleting task", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	s.logger.Info("Task deleted successfully", zap.Int("id", id))
	return nil
}

// package service

// func DeleteTask(w http.ResponseWriter, r *http.Request) {
// 	initlogger()

// 	// Start timer for request duration
// 	start := time.Now()
// 	defer func() {
// 		duration := time.Since(start)
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
// 	}()

// 	params := mux.Vars(r)
// 	id, err := strconv.Atoi(params["id"])
// 	if err != nil {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
// 		http.Error(w, "Invalid ID format", http.StatusBadRequest)
// 		logger.Error("Invalid ID format", zap.Error(err))
// 		return
// 	}
// 	fmt.Printf("Deleting task with ID: %d\n", id)
// 	logger.Info("Deleting task", zap.Int("id", id))

// 	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
// 	if err != nil {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
// 		IncrementDBErrors()
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		logger.Error("Failed to delete task", zap.Error(err))
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
// 		IncrementDBErrors()
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		logger.Error("Failed to get rows affected", zap.Error(err))
// 		return
// 	}
// 	if rowsAffected == 0 {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusNotFound, time.Since(start))
// 		http.Error(w, "Task not found", http.StatusNotFound)
// 		logger.Error("Task not found", zap.Int("id", id))
// 		return
// 	}

// 	IncrementTasksDeleted()

// 	fmt.Printf("Task with ID %d deleted successfully\n", id)
// 	logger.Info("Task deleted successfully", zap.Int("id", id))
// 	// w.WriteHeader(http.StatusNoContent)
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
