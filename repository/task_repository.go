package repository

import (
	"database/sql"
	"time"

	"2do.com/models"
	"go.uber.org/zap"
)

type TaskRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewTaskRepository(db *sql.DB, logger *zap.Logger) *TaskRepository {
	return &TaskRepository{db: db, logger: logger}
}

func (r *TaskRepository) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	rows, err := r.db.Query("SELECT id, title, content, created, status FROM tasks")
	if err != nil {
		r.logger.Error("Failed to fetch tasks", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		var createdBytes []byte

		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &createdBytes, &task.Status); err != nil {
			r.logger.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}

		created, err := time.Parse("2006-01-02 15:04:05", string(createdBytes))
		if err != nil {
			r.logger.Error("Failed to parse datetime", zap.Error(err))
			return nil, err
		}
		task.Created = created

		tasks = append(tasks, task)
		r.logger.Info("Task fetched successfully", zap.Int("id", task.ID))
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(id int) (int64, error) {
	result, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		r.logger.Error("Database error deleting task", zap.Error(err))
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return 0, err
	}
	r.logger.Info("Task deleted Sucessfully", zap.Int("id", id))
	return rowsAffected, nil
}

func (r *TaskRepository) CreateTask(tittle string, content string) (string, error) {
	result, err := r.db.Exec("INSERT INTO tasks (title, content, status, created) VALUES (?, ?, ?, UTC_TIMESTAMP())",
		tittle,
		content,
		"pending",
	)
	if err != nil {
		r.logger.Error("failed to insert task", zap.Error(err))
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("failed to insert ID", zap.Error(err))
		return 0, err
	}
	r.logger.Info("Task created sucessfully", zap.Int("id", id))
	return id, nil
}

func (r *TaskRepository) UpdateTask(id int, tittle string, content string) (int64, error) {
	result, err := r.db.Exec("UPDATE tasks SET title = ?, content = ?, status = ? WHERE id = ?",
		task.Title,
		task.Content,
		task.Status,
		id,
	)
	if err != nil {
		r.logger.Error("failed to update task", zap.Int("id", id), zap.Error(err))
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to update task", zap.Int("id", id), zap.Error(err))
		return 0, err
	}

	r.logger.Info("Task updated successfully", zap.Int("id", id))
	return rowsAffected, nil
}

func (r *TaskRepository) SpecificTask(id int) (int64, error) {
	var task models.Task
	var createdBytes []byte

	err := r.db.Exec("SELECT  id, title, content, created, status FROM tasks WHERE id = ?", id).Scan(
		&task.ID,
		&task.Title,
		&task.Content,
		&createdBytes,
		&task.Status,
	)
	if err == sql.ErrNoRows {
		r.logger.Warn("Task not found", zap.Int("id", id))
		return nil, nil
	}
	if err != nil {
		r.logger.Error("Failed to get task", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	// Parse created time
	created, err := time.Parse("2006-01-02 15:04:05", string(createdBytes))
	if err != nil {
		r.logger.Error("Failed to parse datetime", zap.Int("id", id), zap.Error(err))
		return nil, err
	}
	task.Created = created

	r.logger.Info("Task retrieved successfully",
		zap.Int("id", task.ID))

	return &task, nill
}
