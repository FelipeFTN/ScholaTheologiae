package controllers

import "scholatheologiae-api/models"

func (c *Controllers) Search(query string) ([]models.SearchResult, error) {
	res, err := c.svc.Search(query)
	if err != nil {
		return nil, nil
	}

	return res, nil
}
