package service

import (
	"errors"
	"time"

	"../models"
	"../repository"
	"go.uber.org/zap"
)

type TaskService struct {
	logger *zap.Logger
	repo   *repository.TaskRepository
}

func (s *TaskService) UpdateTask(id int, tittle string, content string, status string) error {
	if id <= 0 {
		return errors.New("invalid task ID")
	}
	if tittle == "" {
		return errors.New("invalid task title")
	}
	if content == "" {
		return errors.New("invalid task content")
	}
	if status == "" {
		return errors.New("invalid task status")
	}

	err := s.repo.UpdateTask(id, tittle, content, status)
	if err != nil {
		s.logger.Error("failed to update task", zap.int("id", id), zap.string("tittle", tittle), zap.Error(err))
		return err
	}
	s.logger.info("Task updated sucessfully", zap.int("id", id), zap.string("tittle", tittle))
	return nil
}

//func updateTask(w http.ResponseWriter, r *http.Request) {
//	initlogger()
//
//	// Start timer for request duration
//	start := time.Now()
//	defer func() {
//		duration := time.Since(start)
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
//	}()
//
//	params := mux.Vars(r)
//	id, err := strconv.Atoi(params["id"])
//	if err != nil {
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//		http.Error(w, "Invalid ID format", http.StatusBadRequest)
//		logger.Error("Invalid ID format", zap.Error(err))
//		return
//	}
//
//	var task Task
//	err = json.NewDecoder(r.Body).Decode(&task)
//	if err != nil {
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		logger.Error("Failed to decode request body", zap.Error(err))
//		return
//	}
//
//	_, err = db.Exec("UPDATE tasks SET title = ?, content = ?, status = ? WHERE id = ?",
//		task.Title, task.Content, task.Status, id)
//	if err != nil {
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//		IncrementDBErrors()
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Error("Failed to update task", zap.Error(err))
//		return
//	}
//
//	IncrementTasksCompleted()
//
//	task.ID = id
//	logger.Info("Task updated successfully", zap.Int("id", task.ID))
//	http.Redirect(w, r, "/", http.StatusSeeOther)
//}
