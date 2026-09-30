package models

import (
	"strings"

	"scholatheologiae-api/constants"
)

type BookRequest struct {
	Name    string
	Part    string
	Chapter string
	// Article is reserved for article-level routes. No book ships articles yet,
	// so the API answers 404 when it is present (see controllers.Read).
	Article string
	Type    constants.RequestType
}

// Validate normalizes the request values and decides which data call the request
// maps to.
func (b *BookRequest) Validate() {
	b.Name = normalize(b.Name)
	b.Part = normalize(b.Part)
	b.Chapter = normalize(b.Chapter)
	b.Article = normalize(b.Article)

	switch {
	case b.Name == "":
		// No book given at all: let the data layer answer "book not found".
		b.Type = constants.ListParts
	case b.Part == "":
		b.Type = constants.ListParts
	case b.Chapter == "":
		b.Type = constants.ListChapters
	case b.Article == "":
		b.Type = constants.GetChapter
	default:
		b.Type = constants.GetArticle
	}
}

func normalize(value string) string {
	if value == "" {
		return ""
	}

	// Normalize the value to lowercase, trim whitespace and use underscores instead of spaces
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(value, " ", "_")))
}
