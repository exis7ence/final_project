package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(255) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX scheduler_date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {
	install := false

	_, err := os.Stat(dbFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("не удалось проверить файл базы данных: %w", err)
		}

		install = true
	}

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	if install {
		if _, err := database.Exec(schema); err != nil {
			database.Close()
			return fmt.Errorf("не удалось создать структуру базы данных: %w", err)
		}
	}

	db = database

	return nil
}