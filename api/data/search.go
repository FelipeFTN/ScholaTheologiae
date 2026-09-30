package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"scholatheologiae-api/models"
)

// maxResultsPerBook caps how many hits a single book contributes to a search.
const maxResultsPerBook = 50

// snippetTokens is how many tokens of context a search snippet carries.
const snippetTokens = 18

// wordRe matches the tokens that are worth searching for.
var wordRe = regexp.MustCompile(`[\p{L}\p{N}]`)

// Search performs an accent-insensitive search across all book databases.
//
// It queries the FTS5 index built by tools/indexer (chapter titles and bodies) and
// falls back to a plain title match when a book has no index yet.
//
// Results are grouped by book in a stable alphabetical order and ordered by
// relevance inside each book, so the same query always returns the same list.
func (d *Data) Search(query string) ([]models.SearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("%w: query must not be empty", models.ErrBadRequest)
	}

	var results []models.SearchResult

	for _, name := range d.SQLite.BookNames() {
		database, err := d.SQLite.book(name)
		if err != nil {
			return nil, err
		}

		var bookResults []models.SearchResult
		if d.SQLite.hasFTS(name) {
			bookResults, err = searchFullText(database.db, name, query)
		} else {
			slog.Warn("Full-text index missing for book, falling back to title search", "book", name)
			bookResults, err = searchTitles(database.db, name, query)
		}
		if err != nil {
			return nil, err
		}

		results = append(results, bookResults...)
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

// searchFullText queries the book's FTS5 index and returns the best matches,
// each with a snippet of the text that matched.
func searchFullText(db *sql.DB, name, query string) ([]models.SearchResult, error) {
	table := name + "_fts"

	statement := `SELECT chapter_title, part_title, chapter_number,
	                     snippet(` + table + `, -1, '<mark>', '</mark>', '…', ?)
	              FROM ` + table + `
	              WHERE ` + table + ` MATCH ?
	              ORDER BY rank, chapter_number
	              LIMIT ?`

	rows, err := db.Query(statement, snippetTokens, ftsQuery(query), maxResultsPerBook)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	seen := make(map[string]bool)
	for rows.Next() {
		var result models.SearchResult
		result.Book = name
		err = rows.Scan(&result.ChapterTitle, &result.PartTitle, &result.ChapterNumber, &result.Snippet)
		if err != nil {
			return nil, err
		}

		// Results are ordered by relevance, so the first hit of a chapter is the
		// one to keep. Guards against an index built before chapters were
		// collapsed to one row per (part, chapter).
		key := result.PartTitle + "\x00" + strconv.Itoa(result.ChapterNumber)
		if seen[key] {
			continue
		}
		seen[key] = true

		results = append(results, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// searchTitles is the fallback used when a book has no full-text index.
func searchTitles(db *sql.DB, name, query string) ([]models.SearchResult, error) {
	statement := `SELECT id, part_title, chapter_title, chapter_number
	              FROM ` + name + `
	              WHERE accent_insensitive_like(chapter_title, ?)
	              GROUP BY part_title, chapter_title, chapter_number
	              ORDER BY chapter_title
	              LIMIT ?`

	rows, err := db.Query(statement, query, maxResultsPerBook)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var result models.SearchResult
		result.Book = name
		err = rows.Scan(&result.ID, &result.PartTitle, &result.ChapterTitle, &result.ChapterNumber)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// ftsQuery turns free text into a safe FTS5 query: every token is quoted (so
// operators typed by the user are matched as literal text). Tokens of four or
// more characters also get a trailing "*" so that "graç" already finds "graças",
// while short words like "de" only match themselves.
func ftsQuery(query string) string {
	tokens := make([]string, 0, 8)

	for _, field := range strings.Fields(query) {
		if !wordRe.MatchString(field) {
			continue
		}

		escaped := strings.ReplaceAll(field, `"`, `""`)
		token := `"` + escaped + `"`
		if len([]rune(field)) >= 4 {
			token += "*"
		}

		tokens = append(tokens, token)
	}

	return strings.Join(tokens, " ")
}
