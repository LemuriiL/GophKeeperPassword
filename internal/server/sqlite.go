package server

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/LemuriiL/GophKeeperPassword/internal/model"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	s := &SQLite{db: db}

	if err = s.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	return s, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) migrate(ctx context.Context) error {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	salt TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	type TEXT NOT NULL,
	title TEXT NOT NULL,
	meta TEXT NOT NULL,
	ciphertext TEXT NOT NULL,
	nonce TEXT NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id)
);
`

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SQLite) CreateUser(ctx context.Context, user model.User) (int64, error) {
	res, err := s.db.ExecContext(
		ctx,
		`INSERT INTO users(login, password_hash, salt, created_at) VALUES(?, ?, ?, ?)`,
		user.Login,
		user.PasswordHash,
		user.Salt,
		user.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *SQLite) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	row := s.db.QueryRowContext(
		ctx,
		`SELECT id, login, password_hash, salt, created_at FROM users WHERE login = ?`,
		login,
	)

	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
	)

	return user, err
}

func (s *SQLite) UpsertItem(ctx context.Context, item model.Item) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO items(id, user_id, type, title, meta, ciphertext, nonce, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
	type = excluded.type,
	title = excluded.title,
	meta = excluded.meta,
	ciphertext = excluded.ciphertext,
	nonce = excluded.nonce,
	updated_at = excluded.updated_at`,
		item.ID,
		item.UserID,
		item.Type,
		item.Title,
		item.Meta,
		item.Ciphertext,
		item.Nonce,
		item.UpdatedAt,
	)

	return err
}

func (s *SQLite) ListItems(ctx context.Context, userID int64) ([]model.Item, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, title, meta, ciphertext, nonce, updated_at
FROM items
WHERE user_id = ?
ORDER BY updated_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Item, 0)

	for rows.Next() {
		var item model.Item

		if err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Type,
			&item.Title,
			&item.Meta,
			&item.Ciphertext,
			&item.Nonce,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (s *SQLite) DeleteItem(ctx context.Context, userID int64, id string) error {
	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM items WHERE user_id = ? AND id = ?`,
		userID,
		id,
	)
	return err
}

func Now() time.Time {
	return time.Now().UTC()
}
