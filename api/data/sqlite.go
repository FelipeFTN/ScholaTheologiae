package data

/*
#include "sqlite.h"
#cgo LDFLAGS: -lsqlite3
*/
import "C"
import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"unsafe"

	"github.com/mattn/go-sqlite3"

	"scholatheologiae-api/models"
)

// SQLite is a struct that holds database connections
type SQLite struct {
	databases map[string]SQLiteDB
}

type SQLiteDB struct {
	db *sql.DB
}

// NewSQLite creates database connections with Portuguese accent-insensitive search
func NewSQLite(dbConnections map[string]string) (*SQLite, error) {
	databases := make(map[string]SQLiteDB)

	for key, path := range dbConnections {
		// Check if the database file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("database file does not exist: %s", path)
		}

		// Register a custom driver with our C-based accent-insensitive search function
		driverName := "sqlite3_with_accent_search_" + key
		sql.Register(driverName, &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// Register our C-based accent_insensitive_like function
				err := conn.RegisterFunc("accent_insensitive_like", callAccentInsensitiveLike, true)
				if err != nil {
					return fmt.Errorf("failed to register accent_insensitive_like function: %w", err)
				}

				return nil
			},
		})

		// Open the database with UTF-8 encoding using our custom driver
		connectionString := path + "?_encoding=UTF8&_loc=auto"
		db, err := sql.Open(driverName, connectionString)
		if err != nil {
			return nil, err
		}

		// Check connectivity and set UTF-8 encoding
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %s", err)
		}

		_, err = db.Exec("PRAGMA encoding = 'UTF-8'")
		if err != nil {
			return nil, fmt.Errorf("failed to set UTF-8 encoding: %s", err)
		}

		slog.Info("Database connection established", "database", key)
		databases[key] = SQLiteDB{db}
	}

	return &SQLite{databases: databases}, nil
}

// BookNames returns the known book names in a stable alphabetical order, so that
// iterating them never depends on Go's randomized map ordering.
func (h *SQLite) BookNames() []string {
	names := make([]string, 0, len(h.databases))
	for name := range h.databases {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// book returns the connection of a known book. The name is interpolated into SQL
// (a table name cannot be a bound parameter), so anything that is not part of the
// known library is rejected here instead of reaching a query.
func (h *SQLite) book(name string) (SQLiteDB, error) {
	database, ok := h.databases[name]
	if !ok || database.db == nil {
		return SQLiteDB{}, fmt.Errorf("%w: no database for book '%s'", models.ErrNotFound, name)
	}

	return database, nil
}

// hasFTS reports whether the full-text index built by tools/indexer exists for a
// book.
func (h *SQLite) hasFTS(name string) bool {
	database, err := h.book(name)
	if err != nil {
		return false
	}

	var count int
	err = database.db.QueryRow(
		`SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?`,
		name+"_fts",
	).Scan(&count)

	return err == nil && count > 0
}

// callAccentInsensitiveLike implements Portuguese accent-insensitive string matching
// It takes a text (like "matrimônio") and a pattern (like "matrimonio")
// and returns true if the pattern is found in the text, ignoring accents
func callAccentInsensitiveLike(text, pattern string) bool {
	// Step 1: Convert Go strings to C strings so we can pass them to C functions
	cText := C.CString(text)
	cPattern := C.CString(pattern)
	// Important: Free the C strings when we're done to avoid memory leaks
	defer C.free(unsafe.Pointer(cText))
	defer C.free(unsafe.Pointer(cPattern))

	// Step 2: Get the lengths of our strings
	textLen := C.int(len(text))
	patternLen := C.int(len(pattern))

	// Step 3: Allocate memory for the normalized (accent-free) versions
	// We need enough space for both strings, so we add their lengths + some extra space
	bufferSize := C.int(len(text) + len(pattern) + 10) // +10 for safety

	// Allocate C memory for the normalized strings
	normalizedText := (*C.char)(C.malloc(C.size_t(bufferSize)))
	normalizedPattern := (*C.char)(C.malloc(C.size_t(bufferSize)))

	// Always free allocated memory when done
	defer C.free(unsafe.Pointer(normalizedText))
	defer C.free(unsafe.Pointer(normalizedPattern))

	// Check if memory allocation failed
	if normalizedText == nil || normalizedPattern == nil {
		return false
	}

	// Step 4: Normalize both strings (remove accents)
	// These variables will hold the actual lengths after normalization
	normalizedTextLen := bufferSize
	normalizedPatternLen := bufferSize

	// Call our C function to normalize "matrimônio" -> "matrimonio"
	C.normalize_string_wrapper(cText, textLen, normalizedText, &normalizedTextLen)

	// Call our C function to normalize "matrimonio" -> "matrimonio" (no change in this case)
	C.normalize_string_wrapper(cPattern, patternLen, normalizedPattern, &normalizedPatternLen)

	// Step 5: Check if the normalized pattern is found in the normalized text
	// This is like asking: does "matrimonio" contain "matrimonio"? -> Yes!
	result := C.contains_pattern_wrapper(normalizedText, normalizedTextLen, normalizedPattern, normalizedPatternLen)

	// Convert C result (0 or 1) to Go boolean
	return result != 0
}

// Close closes all database connections
func (h *SQLite) Close() error {
	for _, d := range h.databases {
		if err := d.db.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err.Error())
			return err
		}
	}
	return nil
}
