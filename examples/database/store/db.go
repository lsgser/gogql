/*
|--------------------------------------------------------------------------
| Database example — SQLite setup
|--------------------------------------------------------------------------
|
| Opens SQLite, creates the users table, and seeds demo rows.
|
*/

package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// OpenSQLite opens a SQLite database and runs migrations/seed.
func OpenSQLite(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	seed := []struct{ id, name, email string }{
		{"1", "Ada Lovelace", "ada@example.com"},
		{"2", "Grace Hopper", "grace@example.com"},
	}
	for _, u := range seed {
		if _, err := db.Exec(`INSERT INTO users (id, name, email) VALUES (?, ?, ?)`, u.id, u.name, u.email); err != nil {
			return fmt.Errorf("seed user %s: %w", u.id, err)
		}
	}
	return nil
}
