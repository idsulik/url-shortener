package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/idsulik/url-shortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = createTables(db)

	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Storage{db: db}, nil
}

func createTables(db *sql.DB) error {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
	`)

	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(fmt.Sprintf("failed to close statement: %s", err.Error()))
		}
	}(stmt)

	_, err = stmt.Exec()

	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (s *Storage) SaveUrl(alias, url string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO urls (alias, url) VALUES (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(fmt.Sprintf("failed to close statement: %s", err.Error()))
		}
	}(stmt)

	res, err := stmt.Exec(alias, url)

	if err != nil {
		if sqliteError, ok := err.(sqlite3.Error); ok && sqliteError.Code == sqlite3.ErrConstraint {
			err = storage.ErrorUrlExists
		}

		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")

	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(fmt.Sprintf("failed to close statement: %s", err.Error()))
		}
	}(stmt)

	var url string

	err = stmt.QueryRow(alias).Scan(&url)

	if err != nil {
		if err == sql.ErrNoRows {
			err = storage.ErrUrlNotFound
		}

		return "", fmt.Errorf("failed to execute statement: %w", err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")

	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(fmt.Sprintf("failed to close statement: %s", err.Error()))
		}
	}(stmt)

	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
