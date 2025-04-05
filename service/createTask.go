package service

import (
	"../repository"
	"Errors"
	"go.uber.org/zap"
)

type TaskService struct {
	repo   *repository.TaskRepository
	logger *zap.Logger
}

func NewTaskService(repo *repository.TaskRepository, logger *zap.Logger) *TaskService {
	return &TaskService{repo: repo, logger: logger}
}

func (s *TaskService) CreateTask(tittle string, content string) (string, error) {
	if tittle == "" {
		return "", errors.New("Tittle cannot be empty")
	}
	if len(tittle) > 100 {
		return "", errors.New("Title too long")
	}
	if content == "" {
		return "", errors.New("Content cannot be empty")
	}

	id, err := s.repo.CreateTask(tittle, content)
	if err != nil {
		s.logger.Error("Failed to create task", zap.string("tittle", tittle), zap.Error(err))
		return "", err
	}
	s.logger.infor("Task created", id)
	return id, nil

}

// package service

// import (
// 	"net/http"
// 	"go.uber.org/zap"

// )

// func initlogger() {
// 	var err error
// 	logger, _ = zap.NewProduction()
// 	if err != nil {
// 		log.Fatal("Failed to initialize logger", zap.Error(err))
// 	}
// 	defer logger.Sync()
// }

// func createTask(w http.ResponseWriter, r *http.Request) {
// 	initlogger()

// 	// Start timer for request duration
// 	start := time.Now()
// 	defer func() {
// 		duration := time.Since(start)
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusOK, duration)
// 	}()

// 	err := r.ParseForm()
// 	if err != nil {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		logger.Error("Failed to parse form", zap.Error(err))
// 		return
// 	}

// 	// Extract form values
// 	title := r.FormValue("title")
// 	content := r.FormValue("content")

// 	// Insert the task into the database
// 	result, err := db.Exec("INSERT INTO tasks (title, content, status, created) VALUES (?, ?, ?, UTC_TIMESTAMP())",
// 		title, content, "pending") // Default status is "pending"
// 	if err != nil {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
// 		IncrementDBErrors()
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		logger.Error("Failed to insert task", zap.Error(err))
// 		return
// 	}

// 	IncrementTasksCreated()
// 	// Get the ID of the newly inserted task
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		InstrumentHTTPRequest(r.Method, r.URL.Path, http.StatusBadRequest, time.Since(start))
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		logger.Error("Failed to get last insert ID", zap.Error(err))
// 		return
// 	}

// 	fmt.Printf("New task created with ID: %d\n", id)
// 	logger.Info("New task created", zap.Int64("id", id))

// 	// Redirect to the task list page
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
