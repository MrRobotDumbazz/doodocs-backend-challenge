package repository

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func New(repositoryPath string) (*Repository, error) {
	const op = "repository.sqlite.NewRepository"
	db, err := sql.Open("sqlite3", repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	smpt, err := db.Prepare(
		`
		CREATE TABLE IF NOT EXIST "smtp" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
		"email" TEXT UNIQUE NOT NULL,
		"password" TEXT NOT NULL
		);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = smpt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Repository{db: db}, nil
}
