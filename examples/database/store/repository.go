package store

import (
	"context"
	"database/sql"
	"errors"
)

// UserRecord is a row from the users table.
type UserRecord struct {
	ID    string
	Name  string
	Email string
}

// UserRepository loads users from SQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a repository backed by db.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ByID(ctx context.Context, id string) (*UserRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id = ?`, id)
	var u UserRecord
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]UserRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserRecord
	for rows.Next() {
		var u UserRecord
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
