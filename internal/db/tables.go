package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TablesRepo struct {
	db *sqlx.DB
}

func NewTablesRepo(db *sqlx.DB) TablesRepo {
	return TablesRepo{db: db}
}

func (r TablesRepo) GetAllDBTables() ([]string, error) {
	var tables []string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';`
	if err := r.db.Select(&tables, query); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	return tables, nil
}
