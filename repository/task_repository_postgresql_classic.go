package repository

import (
	"2do.com/models"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgconn"
)

type PostgresqlClassic struct {
	db *sql.DB
}

func NewPostgresqlClassic(db *sql.DB) *PostgresqlClassic {
	return &PostgresqlClassic{
		db: db,
	}
}

// separate the migrate function

func (r *PostgresqlClassic) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        content VARCHAR(255) NOT NULL,
        completed BOOLEAN NOT NULL DEFAULT FALSE,
        created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        priority VARCHAR(50) DEFAULT 'medium'
    ) 
    `
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresqlClassic) CreateTask(ctx context.Context, task models.Task) (*models.Task, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO tasks (title, content, completed, created, priority) VALUES ($1, $2, $3, $4, $5) RETURNING id", task.Title, task.Content, task.Completed, task.Created, task.Priority).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	task.ID = int(id)
	return &task, nil
}

func (r *PostgresqlClassic) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	log.Println("Repository GetAllTasks called")
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, content, completed, created, priority FROM tasks")
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var GetAllTasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.Completed, &task.Created, &task.Priority); err != nil {
			log.Printf("Error scanning task: %v", err)
			return nil, err
		}
		GetAllTasks = append(GetAllTasks, task)
	}
	log.Printf("Repository retrieved %d tasks", len(GetAllTasks))
	return GetAllTasks, nil
}

func (r *PostgresqlClassic) GetbyID(ctx context.Context, id int) (models.Task, error) {
	var task models.Task
	err := r.db.QueryRowContext(ctx, "SELECT id, title, content, completed, created, priority FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Content, &task.Completed, &task.Created, &task.Priority)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (r *PostgresqlClassic) UpdateTask(ctx context.Context, id int, updated models.Task) (*models.Task, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE tasks SET title = $1, content = $2, completed = $3, priority = $4 WHERE id = $5", updated.Title, updated.Content, updated.Completed, updated.Priority, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &updated, nil
}

func (r *PostgresqlClassic) DeleteTask(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
