package handlers

import (
	"2do.com/models"
	"2do.com/repository"
	"2do.com/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type taskHandler struct {
	TaskRepo  repository.TaskRepository
	Templates *template.Template
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
	return &taskHandler{
		TaskRepo:  repo,
		Templates: templates,
	}
}

func (h *taskHandler) Home(w http.ResponseWriter, r *http.Request) {
	tasks, err := service.GetAllTasks(r.Context(), h.TaskRepo)
	if err != nil {
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
	log.Println("GetAllTasks handler called")
	tasks, err := service.GetAllTasks(r.Context(), h.TaskRepo)
	if err != nil {
		log.Printf("Error getting all tasks: %v", err)
		http.Error(w, "Error getting all tasks", http.StatusInternalServerError)
		return
	}
	log.Printf("Number of tasks retrieved: %d", len(tasks))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *taskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	task, err := service.GetTaskByID(r.Context(), h.TaskRepo, id)
	if err != nil {
		if err == repository.ErrNotExist {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task data: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := service.CreateTask(r.Context(), h.TaskRepo, task)
	if err != nil {
		http.Error(w, "Error creating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task data: "+err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = id
	updatedTask, err := service.UpdateTask(r.Context(), h.TaskRepo, task)
	if err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = service.DeleteTask(r.Context(), h.TaskRepo, id)
	if err != nil {
		if err == repository.ErrDeleteFailed {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
