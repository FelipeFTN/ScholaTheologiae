package controllers

import "scholatheologiae-api/models"

func (c *Controllers) Search(query string) ([]models.SearchResult, error) {
	return c.svc.Search(query)
}
