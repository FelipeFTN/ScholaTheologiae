package data

import "fmt"

func (d *Data) GetBooks() ([]string, error) {
	// Prepare the statement
	query := "SELECT name FROM master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY id"
	stmt, err := d.SQLite.databases["master"].db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a slice
	var books []string
	for rows.Next() {
		var book string
		err = rows.Scan(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (d *Data) GetBookParts(book_name string) ([]string, error) {
	// Prepare the statement
	query := fmt.Sprintf("SELECT DISTINCT part_title FROM %s ORDER BY id", book_name)

	// Check if the database for the book exists
	database := d.SQLite.databases[book_name]
	if database.db == nil {
		return nil, fmt.Errorf("database for book '%s' not found", book_name)
	}

	stmt, err := database.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query()
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

	return parts, nil
}

func (d *Data) GetBookChapters(book_name, part string) (map[string]string, error) {
	// Prepare the statement
	query := fmt.Sprintf("SELECT DISTINCT chapter_number, chapter_title FROM %s WHERE part_title = ?", book_name)
	stmt, err := d.SQLite.databases[book_name].db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(part)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a map
	chapters := make(map[string]string)
	for rows.Next() {
		var chapterNumber, chapterTitle string
		err = rows.Scan(&chapterNumber, &chapterTitle)
		if err != nil {
			return nil, err
		}
		chapters[chapterNumber] = chapterTitle
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(chapters) == 0 {
		return nil, fmt.Errorf("no chapters found for part '%s' in book '%s '", part, book_name)
	}

	return chapters, nil
}
