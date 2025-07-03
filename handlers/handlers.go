package handlers

import (
	"2do.com/middleware"
	"2do.com/models"
	"2do.com/repository"
	"2do.com/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strconv"
)

type taskHandler struct {
	Templates *template.Template
	svc       *service.TaskService
	TaskRepo  repository.TaskRepository
}

//type Response struct {
//	Success bool        `json:"success"`
//	Data    interface{} `json:"data,omitempty"`
//	Error   string      `json:"error,omitempty"`
//}

//func NewTaskHandler(repo repository.TaskRepository) *taskHandler {
//	return &taskHandler{TaskRepo: repo}
//}

func NewTaskHandler(repo repository.TaskRepository) *taskHandler {
	templates := template.Must(template.ParseFiles("./ui/html/index.html"))
	logger := middleware.ZapLogger
	svc := service.NewTaskService(repo, logger)
	return &taskHandler{
		svc:       svc,
		Templates: templates,
		TaskRepo:  repo,
	}
}

func (h *taskHandler) Home(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAllTasks(r.Context(), h.TaskRepo)
	logger := middleware.ZapLogger
	if err != nil {
		logger.Error("Error getting all tasks", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Todos []models.Task
	}{
		Todos: tasks,
	}

	h.Templates.ExecuteTemplate(w, "index.html", data)
}

func (h *taskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	logger := middleware.ZapLogger
	tasks, err := h.svc.GetAllTasks(r.Context(), h.TaskRepo)
	if err != nil {
		logger.Error("Error getting all tasks", zap.Error(err))
		http.Error(w, "Error getting all tasks", http.StatusInternalServerError)
		return
	}
	logger.Info("Tasks retrieved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *taskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger := middleware.ZapLogger
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Error("Error getting task by ID", zap.Int("task_id", id))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	task, err := h.svc.GetTaskByID(r.Context(), h.TaskRepo, id)
	if err != nil {
		if err == repository.ErrNotExist {
			logger.Error("Task not found", zap.Int("task_id", id))
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			logger.Error("Error retrieving task", zap.Int("task_id", id), zap.Error(err))
			http.Error(w, "Error retrieving task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	logger := middleware.ZapLogger
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logger.Error("Error decoding task data", zap.Error(err))
		http.Error(w, "Invalid task data: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := h.svc.CreateTask(r.Context(), h.TaskRepo, task)
	if err != nil {
		logger.Error("Error creating task", zap.Error(err))
		http.Error(w, "Error creating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger := middleware.ZapLogger
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Error("Error updating task", zap.Int("task_id", id))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logger.Error("Error decoding task data", zap.Error(err))
		http.Error(w, "Invalid task data: "+err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = id
	updatedTask, err := h.svc.UpdateTask(r.Context(), h.TaskRepo, task)
	if err != nil {
		logger.Error("Error updating task", zap.Int("task_id", id), zap.Error(err))
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger := middleware.ZapLogger
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Error("Error deleting task", zap.Int("task_id", id))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.svc.DeleteTask(r.Context(), h.TaskRepo, id)
	if err != nil {
		if err == repository.ErrDeleteFailed {
			logger.Error("Error deleting task", zap.Int("task_id", id), zap.Error(err))
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			logger.Error("Error deleting task", zap.Int("task_id", id), zap.Error(err))
			http.Error(w, "Error deleting task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
