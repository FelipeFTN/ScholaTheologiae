package data

import (
	"log/slog"
	"scholatheologiae-api/models"
)

// Search performs an accent-insensitive search across all book databases
func (d *Data) Search(query string) ([]models.SearchResult, error) {
	var results []models.SearchResult

	for book_name, database := range d.SQLite.databases {
		if database.db == nil {
			slog.Error("Database for book not found", "book", book_name)
			continue
		}

		// Search using our custom accent-insensitive function
		stmt, err := database.db.Prepare(`
			SELECT id, part_title, chapter_title, chapter_number
			FROM ` + book_name + ` 
			WHERE accent_insensitive_like(chapter_title, ?)
			GROUP BY part_title, chapter_title, chapter_number
			ORDER BY chapter_title
			LIMIT 50
		`)

		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		// Execute the search
		rows, err := stmt.Query(query)
		if err != nil {
			slog.Error("Query execution failed", "error", err.Error())
			return nil, err
		}
		defer rows.Close()

		// Collect results for this database
		var dbResults []models.SearchResult
		for rows.Next() {
			var result models.SearchResult
			result.Book = book_name
			err = rows.Scan(&result.ID, &result.PartTitle, &result.ChapterTitle, &result.ChapterNumber)
			if err != nil {
				return nil, err
			}
			dbResults = append(dbResults, result)
		}

		results = append(results, dbResults...)
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
