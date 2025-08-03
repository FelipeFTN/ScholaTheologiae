package models

import (
	"scholatheologiae-api/constants"
	"strings"
)

type BookRequest struct {
	Name    string
	Part    string
	Chapter string
	Article string // not used yet
	Type    string
}

func (b *BookRequest) Validate() {
	b.Name = normalize(b.Name)
	b.Part = normalize(b.Part)
	b.Chapter = normalize(b.Chapter)
	b.Article = normalize(b.Article)

	if b.Name == "" {
		b.Type = constants.TYPE_LIST_BOOKS
		return
	}

	if b.Part == "" {
		b.Type = constants.TYPE_LIST_PARTS
		return
	}

	if b.Chapter == "" {
		b.Type = constants.TYPE_LIST_CHAPTERS
		return
	}

	if b.Article == "" {
		b.Type = constants.TYPE_GET_CHAPTER
		return
	}
}

func normalize(value string) string {
	if value == "" {
		return ""
	}

	// Normalize the value to lowercase, trim whitespace and use underscores instead of spaces
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(value, " ", "_")))
}
