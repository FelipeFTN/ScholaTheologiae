package models

type SearchResult struct {
	ID            int    `json:"id,omitempty"`
	Book          string `json:"book"`
	ChapterTitle  string `json:"chapter_title"`
	ChapterNumber int    `json:"chapter_number"`
	ArticleTitle  string `json:"article,omitempty"`
	ArticleNumber int    `json:"article_number,omitempty"`
	PartTitle     string `json:"part_title,omitempty"`
	// Snippet carries the matching excerpt with <mark> around the matched terms.
	Snippet string `json:"snippet,omitempty"`
}
