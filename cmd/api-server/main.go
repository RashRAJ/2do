package main

import (
	"database/sql"
	"log"
	"net/http"
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
	logger.Info("Server started on port", zap.String("port", "8080"))
	http.ListenAndServe(":8080", router)
}
