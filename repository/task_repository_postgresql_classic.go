package repository

import (
	"2do.com/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresqlClassic struct {
	*BaseRepository // Embed base repository
}

func NewTaskRepository(db *pgxpool.Pool, logger *zap.Logger) *PostgresqlClassic {
	return &PostgresqlClassic{
		BaseRepository: &BaseRepository{
			db:     db,
			logger: logger,
		},
	}
}

func (r *PostgresqlClassic) CreateTask(ctx context.Context, task models.Task) (*models.Task, error) {
	start := time.Now()
	r.logger.Info("creating task",
		zap.String("title", task.Title),
		zap.String("priority", task.Priority),
		zap.Bool("completed", task.Completed),
	)

	var id int64
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Error("failed to begin transaction",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, context.Canceled) {
			r.logger.Warn("failed to rollback transaction", zap.Error(err))
		}
	}()

	err = tx.QueryRow(ctx,
		"INSERT INTO tasks (user_id, title, content, completed, created, priority) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		task.Title, task.Content, task.Completed, task.Created, task.Priority).Scan(&id)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				r.logger.Warn("duplicate task creation attempted",
					zap.String("title", task.Title),
					zap.String("postgres_code", pgxError.Code),
					zap.Duration("duration", time.Since(start)),
				)
				return nil, ErrDuplicate
			}
		}
		r.logger.Error("failed to insert task",
			zap.Error(err),
			zap.String("title", task.Title),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		r.logger.Error("failed to commit transaction",
			zap.Error(err),
			zap.Int64("task_id", id),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}

	task.ID = int(id)
	r.logger.Info("task created successfully",
		zap.Int("task_id", task.ID),
		zap.String("title", task.Title),
		zap.Duration("duration", time.Since(start)),
	)

	return &task, nil
}

func (r *PostgresqlClassic) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	start := time.Now()
	r.logger.Info("retrieving all tasks")

	rows, err := r.db.Query(ctx, "SELECT id, title, content, completed, created, priority FROM tasks")
	if err != nil {
		r.logger.Error("failed to query tasks",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.Completed, &task.Created, &task.Priority); err != nil {
			r.logger.Error("failed to scan task row",
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error occurred during rows iteration",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}

	r.logger.Info("tasks retrieved successfully",
		zap.Int("count", len(tasks)),
		zap.Duration("duration", time.Since(start)),
	)

	return tasks, nil
}

func (r *PostgresqlClassic) GetbyID(ctx context.Context, id int) (models.Task, error) {
	start := time.Now()
	r.logger.Info("retrieving task by ID", zap.Int("task_id", id))

	var task models.Task
	err := r.db.QueryRow(ctx,
		"SELECT id, title, content, completed, created, priority FROM tasks WHERE id = $1",
		id).Scan(&task.ID, &task.Title, &task.Content, &task.Completed, &task.Created, &task.Priority)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Warn("query timeout while retrieving task",
				zap.Int("task_id", id),
				zap.Duration("duration", time.Since(start)),
			)
		} else {
			r.logger.Error("failed to retrieve task",
				zap.Error(err),
				zap.Int("task_id", id),
				zap.Duration("duration", time.Since(start)),
			)
		}
		return models.Task{}, err
	}

	r.logger.Info("task retrieved successfully",
		zap.Int("task_id", id),
		zap.String("title", task.Title),
		zap.Duration("duration", time.Since(start)),
	)

	return task, nil
}

func (r *PostgresqlClassic) UpdateTask(ctx context.Context, id int, updated models.Task) (*models.Task, error) {
	start := time.Now()
	r.logger.Info("updating task",
		zap.Int("task_id", id),
		zap.String("title", updated.Title),
		zap.Bool("completed", updated.Completed),
	)

	commandTag, err := r.db.Exec(ctx,
		"UPDATE tasks SET title = $1, content = $2, completed = $3, priority = $4 WHERE id = $5",
		updated.Title, updated.Content, updated.Completed, updated.Priority, id)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				r.logger.Warn("duplicate constraint violation during update",
					zap.Int("task_id", id),
					zap.String("postgres_code", pgxError.Code),
					zap.Duration("duration", time.Since(start)),
				)
				return nil, ErrDuplicate
			}
		}
		r.logger.Error("failed to update task",
			zap.Error(err),
			zap.Int("task_id", id),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, err
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Warn("no rows affected during update - task may not exist",
			zap.Int("task_id", id),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, ErrUpdateFailed
	}

	r.logger.Info("task updated successfully",
		zap.Int("task_id", id),
		zap.Int64("rows_affected", rowsAffected),
		zap.Duration("duration", time.Since(start)),
	)

	return &updated, nil
}

func (r *PostgresqlClassic) DeleteTask(ctx context.Context, id int) error {
	start := time.Now()
	r.logger.Info("deleting task", zap.Int("task_id", id))

	commandTag, err := r.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		r.logger.Error("failed to delete task",
			zap.Error(err),
			zap.Int("task_id", id),
			zap.Duration("duration", time.Since(start)),
		)
		return err
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Warn("no rows affected during delete - task may not exist",
			zap.Int("task_id", id),
			zap.Duration("duration", time.Since(start)),
		)
		return ErrDeleteFailed
	}

	r.logger.Info("task deleted successfully",
		zap.Int("task_id", id),
		zap.Int64("rows_affected", rowsAffected),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

// Close properly closes the logger
func (r *PostgresqlClassic) Close() {
	if r.logger != nil {
		r.logger.Sync()
	}
}
