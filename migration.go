package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigration(db *sql.DB) error {

	//create migrations table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	//get list of already applied migrations
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to quesry migrations: %w", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		appliedMigrations[version] = true
	}

	//read migration files from schema directory
	files, err := filepath.Glob("sql/schema/[0-9]*.sql")
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}
	sort.Strings(files)

	//apply each migration that hasn't been applied yet
	for _, file := range files {
		//extract version from filename
		baseName := filepath.Base(file)
		version := strings.Split(baseName, "_")[0]

		if appliedMigrations[version] {
			fmt.Printf("Migration %s already applied\n", baseName)
			continue
		}

		//read migration file
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		//Parse content to extract just the Up part
		fileContent := string(content)
		upIndex := strings.Index(fileContent, "-- +goose Up")
		downIndex := strings.Index(fileContent, "-- +goose Down")

		var sqlToExecute string
		if upIndex >= 0 && downIndex > upIndex {
			// extract just the up part
			sqlToExecute = fileContent[upIndex+len("-- +goose Up") : downIndex]
		} else {
			// No goose directives, use the whole file
			sqlToExecute = fileContent
		}

		//apply migration
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		//execute migration
		fmt.Printf("Applying migration %s...\n", baseName)
		if _, err := tx.Exec(sqlToExecute); err != nil {
			fmt.Printf("Error executing migration SQL: %v\n", err)
			//check if it's a "relation already exists" error
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Warning: migration %s tried to create something that already exists. Marking as applied anyway.\n", baseName)

				//still record
				if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to record migration %s: %w", file, err)
				}
				if err := tx.Commit(); err != nil {
					return fmt.Errorf("failed to commit transaction: %w", err)
				}

				fmt.Printf("Successfully marked migration %s as applied\n", baseName)
				continue
			} else {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration %s: %w", file, err)
			}
		}

		//record that migrations has been applied
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		fmt.Printf("Successfully applied migration %s\n", baseName)
	}

	return nil
}
