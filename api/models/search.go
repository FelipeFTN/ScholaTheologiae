package models

type SearchResult struct {
	ID            int    `json:"id"`
	Book          string `json:"book"`
	ChapterTitle  string `json:"chapter_title"`
	ChapterNumber int    `json:"chapter_number"`
	ArticleTitle  string `json:"article,omitempty"`
	ArticleNumber int    `json:"article_number,omitempty"`
	Content       string `json:"content,omitempty"`
	PartTitle     string `json:"part_title,omitempty"`
}
