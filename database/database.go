package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)


func InitDB(filePath string) (*sql.DB, error) {



	DSN := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", filePath)
	DB, err := sql.Open("sqlite3", DSN)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения драйвера базы данных: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		DB.Close()
		return nil, fmt.Errorf("Ошибка подключения базы данных: %v", err)
	}

	err = migrate(DB)
	if err != nil {
		return nil, err
	}

	return DB, nil

}

func migrate(db *sql.DB) error {

	createTableSiteSQlite := `CREATE TABLE IF NOT EXISTS sites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL UNIQUE,
    interval_seconds INTEGER NOT NULL DEFAULT 60,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSiteSQlite)
	if err != nil {
		return fmt.Errorf("Ошибка при создании таблицы сайта в базе данных: %v", err)
	}

	createTableCheckSiteSQlite := `CREATE TABLE IF NOT EXISTS check_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id INTEGER NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    is_available BOOLEAN NOT NULL,
    checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(createTableCheckSiteSQlite)
	if err != nil {
		return fmt.Errorf("Ошибка при создании таблицы проверки сайта в базе данных: %v", err)
	}

	return nil

}