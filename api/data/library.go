package data

import (
	"fmt"
	"os"
)

type Library struct {
	db *SQLite
}

func (l *Library) GetBooks() []string {
	files, err := os.ReadDir("./data/library")
	if err != nil {
		return nil
	}

	var books []string
	for _, file := range files {
		if !file.IsDir() {
			books = append(books, file.Name())
		}
	}

	return books
}

func (l *Library) GetChapter(name, part, chapter string) (string, error) {
	filePath := fmt.Sprintf("./data/library/%s/content/%s/chapter_%s.md", name, part, chapter)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
