package handlers

import (
	"2do.com/middleware"
	"net/http"

	"2do.com/repository"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, taskRepo repository.TaskRepository) {
	router.Use(middleware.Logger)

	// Serve static files
	fs := http.FileServer(http.Dir("./ui/html"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.Auth)

	taskHandler := NewTaskHandler(taskRepo)

	router.HandleFunc("/", taskHandler.Home).Methods("GET")
	taskRouter := apiRouter.PathPrefix("/tasks").Subrouter()
	taskRouter.HandleFunc("", taskHandler.GetAllTasks).Methods("GET")
	taskRouter.HandleFunc("", taskHandler.CreateTask).Methods("POST")
	taskRouter.HandleFunc("/{id:[0-9]+}", taskHandler.GetTaskByID).Methods(http.MethodGet)
	taskRouter.HandleFunc("/{id:[0-9]+}", taskHandler.UpdateTask).Methods(http.MethodPut)
	taskRouter.HandleFunc("/{id:[0-9]+}", taskHandler.DeleteTask).Methods(http.MethodDelete)
}
