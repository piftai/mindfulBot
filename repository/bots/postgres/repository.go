package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository - реализация таблицы bots.
type Repository struct {
	db *sqlx.DB
}

// New - конструктор *Repository.
func New(db *sqlx.DB) *Repository {
	return &Repository{db}
}

// Flush - очищает таблицу.
func (r *Repository) Flush() error {
	_, err := r.db.Exec(`TRUNCATE TABLE bots RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("bots.Flush: %w", err)
	}

	return nil
}
