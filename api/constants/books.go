package constants

// RequestType tells the data layer which call a request maps to.
type RequestType int

const (
	// ListParts lists the parts (books, books' sections) of one book.
	ListParts RequestType = iota
	// ListChapters lists the chapters of one part of a book.
	ListChapters
	// GetChapter returns the body of one chapter.
	GetChapter
	// GetArticle is reserved for article-level reads: no book ships articles yet.
	GetArticle
)
