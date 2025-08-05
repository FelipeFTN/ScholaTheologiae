package data

import (
	"log/slog"
	"scholatheologiae-api/models"
)

func (d *Data) Search(query string) ([]models.SearchResult, error) {
	// Prepare the statement for searching across all book databases
	// There is no master yet, so we will search in all databases
	var results []models.SearchResult
	for book_name, database := range d.SQLite.databases {
		if database.db == nil {
			// This should not happen
			slog.Error("Database for book not found", "book", book_name)
			continue
		}

		slog.Info("Searching in book database", "book", book_name, "query", query)
		// Prepare the statement
		stmt, err := database.db.Prepare(`
			SELECT id, part_title, chapter_title, chapter_number
			FROM ` + book_name + ` 
			WHERE LOWER(chapter_title) LIKE LOWER(?)
			GROUP BY id, part_title, chapter_title, chapter_number
			ORDER BY id
			LIMIT 50
		`)

		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		// normalizedQuery := removeAccents(query)
		// Execute the statement
		rows, err := stmt.Query("%" + query + "%")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// Scan the results into a slice
		for rows.Next() {
			var result models.SearchResult
			result.Book = book_name
			err = rows.Scan(&result.ID, &result.PartTitle, &result.ChapterTitle, &result.ChapterNumber)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
