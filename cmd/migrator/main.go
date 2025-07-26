package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sergeyreshetnyakov/notion/internal/config"
)

func main() {
	cfg := config.MustLoad()

	if cfg.StoragePath == "" {
		panic("storage_path is empty")
	}
	if cfg.MigrationsPath == "" {
		panic("migrations_path is empty")
	}

	var migrationsTable string
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.StoragePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("Migrations was successfuly created")
}
