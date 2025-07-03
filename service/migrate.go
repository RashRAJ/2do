package service

import (
	"2do.com/migrate"
	"context"
	"fmt"
)

func Migrate(ctx context.Context, migrator migrate.Migrator) error {
	exists, err := migrator.CheckMigrationStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if !exists {
		if err := migrator.Migrate(ctx); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
