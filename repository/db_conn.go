package repository

import (
	"2do.com/config"
	"2do.com/middleware"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// ConnectDB establishes a database connection pool with retry logic
func ConnectDB() *pgxpool.Pool {
	dbConfig := config.LoadDBConfig()
	logger := middleware.ZapLogger

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.TotalTimeout)
	defer cancel()

	poolConfig, err := buildPoolConfig(dbConfig)
	if err != nil {
		log.Printf("Failed to build pool config: %v", err)
		logger.Error("Failed to build pool config", zap.Error(err))
		return nil
	}

	pool := attemptConnection(ctx, dbConfig, poolConfig)
	if pool == nil {
		return nil
	}

	log.Printf("Database pool connected successfully")
	logger.Info("Database pool connected successfully")
	return pool
}

// buildPoolConfig constructs the pgxpool configuration
func buildPoolConfig(config config.DbConfig) (*pgxpool.Config, error) {
	logger := middleware.ZapLogger

	baseConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName,
	)

	poolConfig, err := pgxpool.ParseConfig(baseConnStr)
	if err != nil {
		logger.Error("failed to parse config", zap.Error(err))
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = int32(config.MaxOpenConns)
	poolConfig.MinConns = int32(config.MaxIdleConns)
	poolConfig.MaxConnLifetime = config.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = config.ConnMaxIdleTime
	poolConfig.HealthCheckPeriod = time.Minute

	return poolConfig, nil
}

// attemptConnection handles the retry logic for database connection pool

func attemptConnection(ctx context.Context, dbConfig config.DbConfig, poolConfig *pgxpool.Config) *pgxpool.Pool {
	logger := middleware.ZapLogger
	backoffStrategy := createBackoffStrategy(dbConfig)

	var pool *pgxpool.Pool
	var lastErr error

	startTime := time.Now()
	retryCount := 0

	for retryCount < dbConfig.MaxRetries {
		// Check if we've exceeded total timeout
		if time.Since(startTime) >= dbConfig.TotalTimeout {
			logger.Error("Failed to connect to database: total timeout exceeded")
			return nil
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			log.Printf("Database connection cancelled: %v", ctx.Err())
			logger.Error("Database connection cancelled", zap.Error(ctx.Err()))
			return nil
		default:
		}

		// Attempt connection
		var err error
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			lastErr = fmt.Errorf("failed to create connection pool: %w", err)
		} else {
			if err := pool.Ping(ctx); err != nil {
				pool.Close()
				lastErr = fmt.Errorf("ping failed: %w", err)
			} else {
				// Success!
				return pool
			}
		}

		retryCount++
		if retryCount >= dbConfig.MaxRetries {
			break
		}

		// Calculate backoff duration
		backoffDuration := backoffStrategy.NextBackOff()
		log.Printf("Database connection failed (attempt %d/%d): %v. Retrying in %v",
			retryCount, dbConfig.MaxRetries, lastErr, backoffDuration)
		logger.Error("Database connection failed", zap.Error(lastErr), zap.Int("retry_count", retryCount))

		// Wait for backoff duration or context cancellation
		select {
		case <-ctx.Done():
			log.Printf("Database connection cancelled during backoff: %v", ctx.Err())
			logger.Error("Database connection cancelled during backoff", zap.Error(ctx.Err()))
			return nil
		case <-time.After(backoffDuration):
			// Continue to next retry
		}
	}

	log.Printf("Failed to connect to database after %d retries: %v", dbConfig.MaxRetries, lastErr)
	logger.Error("Failed to connect to database after retries", zap.Error(lastErr))
	return nil
}

func createBackoffStrategy(config config.DbConfig) *backoff.ExponentialBackOff {
	exponentialBackoff := &backoff.ExponentialBackOff{
		InitialInterval:     config.InitialDelay,
		RandomizationFactor: 0.1,
		Multiplier:          2.0,
		MaxInterval:         config.MaxDelay,
	}

	exponentialBackoff.Reset()
	return exponentialBackoff
}

// ClosePool safely closes the pool
func ClosePool(pool *pgxpool.Pool) {
	logger := middleware.ZapLogger
	if pool != nil {
		pool.Close()
		log.Printf("Database pool closed")
		logger.Info("Database pool closed")
	}
}
