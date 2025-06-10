package db

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Config struct {
	FilePath     string `env:"DB_FILE_PATH,default=sqlite.db"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS,default=2"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS,default=4"`
	LogLevel     string `env:"DB_LOG_LEVEL,default=error"`
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewDB(cfg Config) (*sqlx.DB, error) {
	// Database file path
	dbPath := filepath.Join(cfg.FilePath)

	// Check if the database file exists, create if not
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("create sqlite db file: %w", err)
		}
		file.Close()
	}

	// Connect to SQLite
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("sqlite conn error: %w", err)
	}

	// Set connection pool configurations
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// Test the database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	// Apply migrations
	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("set migrations dialect: %w", err)
	}

	if err = goose.Up(db.Unsafe().DB, "migrations"); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db.Unsafe(), nil
}
