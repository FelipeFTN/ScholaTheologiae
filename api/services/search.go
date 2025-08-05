package services

import "scholatheologiae-api/models"

func (s *Services) Search(query string) ([]models.SearchResult, error) {
	return s.data.Search(query)
}
