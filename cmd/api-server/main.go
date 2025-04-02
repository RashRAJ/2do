package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Status  string    `json:"status"`
}

var db *sql.DB
var logger *zap.Logger

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to my home page"))
}

func initlogger() {
	var err error
	logger, _ = zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger", zap.Error(err))
	}
	defer logger.Sync()
}

func main() {
	initlogger()
	// ExposeMetrics()

	var err error
	db, err = sql.Open("mysql", "web:password@tcp(127.0.0.1:3306)/snippetbox")
	if err != nil {
		logger.Error("Failed to connect to db", zap.Error(err))
	}
	logger.Info("Connected to db")
	defer db.Close()

	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/", allTask).Methods("GET")
	router.HandleFunc("/task/{id}", specificTask).Methods("GET")
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/update/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/delete/{id}", DeleteTask).Methods("DELETE")

	// log.Fatal(http.ListenAndServe(":8080", router))
	logger.Info("Server started on port", zap.String("port", "8080"))
	http.ListenAndServe(":8080", router)
}

func allTask(w http.ResponseWriter, r *http.Request) {
	initlogger()

	// Start timer for request duration
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
	}()

	var tasks []Task
	rows, err := db.Query("SELECT id, title, content, created, status FROM tasks")
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		IncrementDBErrors()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to fetch tasks", zap.Error(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var createdBytes []byte // Temporary variable to hold the raw datetime value

		// Scan the row into the task struct
		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &createdBytes, &task.Status); err != nil {
			InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
			IncrementDBErrors()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Failed to scan row", zap.Error(err))
			return
		}

		// Parse the raw datetime value into a time.Time object
		created, err := time.Parse("2006-01-02 15:04:05", string(createdBytes))
		if err != nil {
			InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Failed to parse datetime", zap.Error(err))
			return
		}
		task.Created = created

		tasks = append(tasks, task)
		logger.Info("Task fetched successfully", zap.Int("id", task.ID))
	}

	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to get executable path", zap.Error(err))
		return
	}
	tmpl, err := template.ParseFiles("ui/html/index.html")
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to parse template", zap.Error(err))
		return
	}

	data := struct {
		Todos []Task
	}{
		Todos: tasks,
	}

	tmpl.Execute(w, data)

}

func specificTask(w http.ResponseWriter, r *http.Request) {
	initlogger()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		logger.Error("Invalid ID formart", zap.Error(err))
		return
	}

	var task Task
	var createdBytes []byte // Temporary variable to hold the raw datetime value

	// Query the database and scan the result
	err = db.QueryRow("SELECT id, title, content, created, status FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.Title, &task.Content, &createdBytes, &task.Status)
	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		logger.Error("Task not found", zap.Error(err))
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to query db", zap.Error(err))
		return
	}

	// Parse the raw datetime value into a time.Time object
	created, err := time.Parse("2006-01-02 15:04:05", string(createdBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to parse datetime", zap.Error(err))
		return
	}
	task.Created = created

	// Encode the task as JSON and send it in the response
	json.NewEncoder(w).Encode(task)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	initlogger()

	// Start timer for request duration
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
	}()

	err := r.ParseForm()
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("Failed to parse form", zap.Error(err))
		return
	}

	// Extract form values
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Insert the task into the database
	result, err := db.Exec("INSERT INTO tasks (title, content, status, created) VALUES (?, ?, ?, UTC_TIMESTAMP())",
		title, content, "pending") // Default status is "pending"
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		IncrementDBErrors()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to insert task", zap.Error(err))
		return
	}

	IncrementTasksCreated()
	// Get the ID of the newly inserted task
	id, err := result.LastInsertId()
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to get last insert ID", zap.Error(err))
		return
	}

	fmt.Printf("New task created with ID: %d\n", id)
	logger.Info("New task created", zap.Int64("id", id))

	// Redirect to the task list page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	initlogger()

	// Start timer for request duration
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
	}()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		logger.Error("Invalid ID format", zap.Error(err))
		return
	}

	var task Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("Failed to decode request body", zap.Error(err))
		return
	}

	_, err = db.Exec("UPDATE tasks SET title = ?, content = ?, status = ? WHERE id = ?",
		task.Title, task.Content, task.Status, id)
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		IncrementDBErrors()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to update task", zap.Error(err))
		return
	}

	IncrementTasksCompleted()

	task.ID = id
	logger.Info("Task updated successfully", zap.Int("id", task.ID))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	initlogger()

	// Start timer for request duration
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
	}()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		logger.Error("Invalid ID format", zap.Error(err))
		return
	}
	fmt.Printf("Deleting task with ID: %d\n", id)
	logger.Info("Deleting task", zap.Int("id", id))

	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		IncrementDBErrors()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to delete task", zap.Error(err))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
		IncrementDBErrors()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to get rows affected", zap.Error(err))
		return
	}
	if rowsAffected == 0 {
		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusNotFound, time.Since(start))
		http.Error(w, "Task not found", http.StatusNotFound)
		logger.Error("Task not found", zap.Int("id", id))
		return
	}

	IncrementTasksDeleted()

	fmt.Printf("Task with ID %d deleted successfully\n", id)
	logger.Info("Task deleted successfully", zap.Int("id", id))
	// w.WriteHeader(http.StatusNoContent)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
