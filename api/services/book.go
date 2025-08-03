package services

func (s *Services) ListParts(name string) ([]string, error) {
	// Get the list of parts from the database
	return s.data.GetBookParts(name)
}

func (s *Services) ListChapters(name, part string) (map[string]string, error) {
	// Get the list of chapters for the specified part from the database
	return s.data.GetBookChapters(name, part)
}

func (s *Services) GetChapter(name, part, chapter string) (string, error) {
	// Fetch the chapter content from the library
	return s.data.Library.GetChapter(name, part, chapter)
}
