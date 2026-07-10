package db

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations создаёт таблицу schema_migrations и применяет все .sql файлы
// из папки migrationsDir в алфавитном порядке, пропуская уже применённые.
func RunMigrations(ctx context.Context, db *pgxpool.Pool, migrationsDir string) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, f := range files {
		name := filepath.Base(f)

		var exists bool
		_ = db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename=$1)`, name,
		).Scan(&exists)
		if exists {
			continue
		}

		sql, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		// Вырезаем goose-комментарии, берём только блок Up
		content := extractUpBlock(string(sql))

		if _, err := db.Exec(ctx, content); err != nil {
			return err
		}

		_, _ = db.Exec(ctx,
			`INSERT INTO schema_migrations (filename) VALUES ($1)`, name,
		)
		slog.Info("migration applied", "file", name)
	}
	return nil
}

// extractUpBlock берёт SQL между "-- +goose Up" и "-- +goose Down" (или конец файла).
func extractUpBlock(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inUp := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "-- +goose Up" {
			inUp = true
			continue
		}
		if trimmed == "-- +goose Down" {
			break
		}
		if inUp {
			result = append(result, line)
		}
	}
	if len(result) == 0 {
		return content // нет goose-маркеров — выполняем весь файл
	}
	return strings.Join(result, "\n")
}
