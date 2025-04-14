package service

import (
	"time"

	"../models"
	"../repository"
	"go.uber.org/zap"
)

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		s.logger.Info("GetAllTasks duration", zap.Duration("duration", duration))
	}()

	tasks, err := s.repo.GetAllTasks()
	if err != nil {
		s.logger.Error("Failed to get tasks", zap.Error(err))
		return nil, err
	}

	return tasks, nil
}

//package service
//
//import (
//	"database/sql"
//	"html/template"
//	"net/http"
//	"time"
//
//	"go.uber.org/zap"
//)
//
//var db *sql.DB
//var logger *zap.Logger
//
//func allTask(w http.ResponseWriter, r *http.Request) {
//	initlogger()
//
//	// Start timer for request duration
//	start := time.Now()
//	defer func() {
//		duration := time.Since(start)
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
//	}()
//
//	var tasks []Task
//	rows, err := db.Query("SELECT id, title, content, created, status FROM tasks")
//	if err != nil {
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//		IncrementDBErrors()
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Error("Failed to fetch tasks", zap.Error(err))
//		return
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var task Task
//		var createdBytes []byte // Temporary variable to hold the raw datetime value
//
//		// Scan the row into the task struct
//		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &createdBytes, &task.Status); err != nil {
//			InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//			IncrementDBErrors()
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			logger.Error("Failed to scan row", zap.Error(err))
//			return
//		}
//
//		// Parse the raw datetime value into a time.Time object
//		created, err := time.Parse("2006-01-02 15:04:05", string(createdBytes))
//		if err != nil {
//			InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			logger.Error("Failed to parse datetime", zap.Error(err))
//			return
//		}
//		task.Created = created
//
//		tasks = append(tasks, task)
//		logger.Info("Task fetched successfully", zap.Int("id", task.ID))
//	}
//
//	if err != nil {
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Error("Failed to get executable path", zap.Error(err))
//		return
//	}
//	tmpl, err := template.ParseFiles("ui/html/index.html")
//	if err != nil {
//		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Error("Failed to parse template", zap.Error(err))
//		return
//	}
//
//	data := struct {
//		Todos []Task
//	}{
//		Todos: tasks,
//	}
//
//	tmpl.Execute(w, data)
//
//}
