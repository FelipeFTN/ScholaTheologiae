package models

type SearchResult struct {
	ID            int    `json:"id"`
	Book          string `json:"book"`
	ChapterTitle  string `json:"chapter_title"`
	ChapterNumber int    `json:"chapter_number"`
	Article       string `json:"article,omitempty"`
	Content       string `json:"content,omitempty"`
	PartTitle     string `json:"part_title,omitempty"`
	PartSubtitle  string `json:"part_subtitle,omitempty"`
}
