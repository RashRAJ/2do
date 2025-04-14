package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"../models"
	"../service"
	"go.uber.org/zap"
)

type TaskHandler struct {
	service *service.TaskService
	logger  *zap.Logger
}

func NewTaskHandler(service *service.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{service: service, logger: logger}
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		h.logger.Error("Failed to get tasks", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("ui/html/index.html")
	if err != nil {
		h.logger.Error("Failed to parse template", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Todos []models.Task
	}{
		Todos: tasks,
	}

	tmpl.Execute(w, data)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error("Invalid ID format", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		h.logger.Error("Failed to delete task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) == 0 || len(content) == 0 {
		h.logger.Error("Invalid title/content format")
		http.Error(w, "Invalid title/content format", http.StatusBadRequest)
		return
	}

	_, err := h.service.CreateTask(title, content)
	if err != nil {
		h.logger.Error("Failed to create task",
			zap.String("title", title),
			zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("failed to parse form", zap.Error(err))
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) == 0 || len(content) == 0 {
		h.logger.Error("Invalid title/content format")
		http.Error(w, "Invalid title/content format", http.StatusBadRequest)
		return
	}
}

func (h *TaskHandler) SpecificTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error("Invalid ID format", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	task, err := h.service.SpecificTask(id)
	if err != nil {
		h.logger.Error("Failed to fetch task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/task/{id}", http.StatusSeeOther)

	if task == nil {
		h.logger.Warn("Task not found",
			zap.Int("id", id))
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
