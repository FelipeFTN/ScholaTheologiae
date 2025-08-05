package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite is a struct that holds the database connection
type SQLite struct {
	// db is the SQLite database connection
	databases map[string]SQLiteDB
}

type SQLiteDB struct {
	db *sql.DB
}

// NewSQLite creates a new SQLite
func NewSQLite(dbConnections map[string]string) (*SQLite, error) {
	// Create a map to hold the database connections
	databases := make(map[string]SQLiteDB)
	for key, path := range dbConnections {
		// Check if the database file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("database file does not exist: %s", path)
		}

		// Open the SQLite database with UTF-8 encoding
		connectionString := path + "?_encoding=UTF8&_loc=auto"
		db, err := sql.Open("sqlite3", connectionString)
		if err != nil {
			return nil, err
		}
		// Check if the database is reachable
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %s", err)
		}

		// Set UTF-8 encoding pragma
		_, err = db.Exec("PRAGMA encoding = 'UTF-8'")
		if err != nil {
			return nil, fmt.Errorf("failed to set UTF-8 encoding: %s", err)
		}

		// Insert the database connection into the map with the key name
		databases[key] = SQLiteDB{db}
		// databases[dbPath] = SQLiteDB{db}
	}

	// return the SQLite
	return &SQLite{
		databases: databases,
	}, nil
}

// Close closes the database connection
func (h *SQLite) Close() error {
	// Iterate over the database connections and close them
	var err error
	for _, d := range h.databases {
		err = d.db.Close()
		if err != nil {
			slog.Error("Failed to close database connection", "error", err.Error())
		}
	}

	return nil
}

// Execute executes a query on the database
func (h *SQLiteDB) Execute(query string, args ...any) (sql.Result, error) {
	// Prepare the statement
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Query executes a query on the database and returns the rows
func (h *SQLiteDB) Query(query string, args ...any) (*sql.Rows, error) {
	// Prepare the statement
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// QueryRow executes a query on the database and returns a single row
func (h *SQLiteDB) QueryRow(query string, args ...any) *sql.Row {
	// Prepare the statement
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return nil
	}
	defer stmt.Close()

	// Execute the statement
	row := stmt.QueryRow(args...)

	return row
}
