package repository

import (
	"context"
	"database/sql"
	"errors"

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
func (r *PostgresqlClassic) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        content VARCHAR(255) NOT NULL,
        completed BOOLEAN NOT NULL DEFAULT FALSE,
        created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    ) 
    `
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresqlClassic) CreateTask(ctx context.Context, task Task) (*Task, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "INSERT INTO tasks (title, content, completed, created) VALUES ($1, $2, $3, $4) RETURNING id", task.Title, task.Content, task.Completed, task.Created).Scan(&id)
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

func (r *PostgresqlClassic) GetAllTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, content, completed, created FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var GetAllTasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.Completed, &task.Created); err != nil {
			return nil, err
		}
		GetAllTasks = append(GetAllTasks, task)
	}
	return GetAllTasks, nil
}

func (r *PostgresqlClassic) GetbyID(ctx context.Context, id int) (Task, error) {
	var task Task
	err := r.db.QueryRowContext(ctx, "SELECT id, title, content, completed, created FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Content, &task.Completed, &task.Created)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *PostgresqlClassic) UpdateTask(ctx context.Context, id int, updated Task) (*Task, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE tasks SET title = $1, content = $2, completed = $3 WHERE id = $4", updated.Title, updated.Content, updated.Completed, id)
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
