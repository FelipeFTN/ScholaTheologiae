package data

import (
	"fmt"
	"strconv"

	"scholatheologiae-api/models"
)

// GetBookParts lists the parts of a book.
func (d *Data) GetBookParts(bookName string) ([]string, error) {
	database, err := d.SQLite.book(bookName)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT DISTINCT part_title FROM %s ORDER BY id", bookName)

	rows, err := database.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a slice
	var parts []string
	for rows.Next() {
		var part string
		err = rows.Scan(&part)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}

// GetBookChapters lists the chapters of one part as "chapter_number: title" pairs.
func (d *Data) GetBookChapters(bookName, part string) (map[string]string, error) {
	database, err := d.SQLite.book(bookName)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT DISTINCT chapter_number, chapter_title FROM %s WHERE part_title = ? ORDER BY chapter_number, id",
		bookName,
	)

	rows, err := database.db.Query(query, part)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a map
	chapters := make(map[string]string)
	for rows.Next() {
		var chapterNumber int
		var chapterTitle string
		err = rows.Scan(&chapterNumber, &chapterTitle)
		if err != nil {
			return nil, err
		}
		chapters[strconv.Itoa(chapterNumber)] = chapterTitle
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(chapters) == 0 {
		return nil, fmt.Errorf("%w: no chapters found for part '%s' in book '%s'", models.ErrNotFound, part, bookName)
	}

	return chapters, nil
}
