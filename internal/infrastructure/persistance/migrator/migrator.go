package migrator

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"

	//nolint:revieve // This is for migrate
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

const (
	requiredPartsLen = 2
)

// RunMigrations runs database migrations using the provided connection string and migrations path.
//
// It prepares the connection string, checks if the migrations path is provided,
// and initiates the migration process. If the context is canceled, it stops the
// migration process gracefully.
//
// Parameters:
//   - ctx: The context to manage cancellation.
//   - connStr: The database connection string.
//   - migrationsPath: The path to the migration files.
//
// Returns:
//   - An error if the migrations fail; otherwise, nil.
func RunMigrations(ctx context.Context, connStr, migrationsPath string) error {
	connStr = prepareDSN(connStr)

	if migrationsPath == "" {
		slog.Debug("Skipping migrations")
		return nil
	}

	slog.Debug("Running migrations")

	migrator, err := migrate.New(
		"file://"+migrationsPath,
		connStr,
	)
	if err != nil {
		slog.Error("Error while creating migrator", slog.Any("err", err))
		return err
	}
	defer migrator.Close()

	errCh := make(chan error)

	go func() {
		errCh <- migrator.Up()
	}()

	select {
	case err = <-errCh:
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("Error while running migrations", slog.Any("err", err))
			return err
		}
	case <-ctx.Done():
		migrator.GracefulStop <- true
		return ctx.Err()
	}

	slog.Debug("Migrations finished")

	return nil
}

func prepareDSN(connStr string) string {
	if strings.HasPrefix(connStr, "pgx5://") {
		return connStr
	}

	if strings.Contains(connStr, "://") {
		parts := strings.Split(connStr, "://")
		if len(parts) == requiredPartsLen {
			return "pgx5://" + parts[1]
		}
		return ""
	}

	if connStr != "" {
		return "pgx5://" + connStr
	}

	return connStr
}
