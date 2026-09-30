package controllers

import (
	"fmt"

	"scholatheologiae-api/constants"
	"scholatheologiae-api/models"
)

func (c *Controllers) Read(request models.BookRequest) (any, error) {
	switch request.Type {
	case constants.ListParts:
		return c.svc.ListParts(request.Name)
	case constants.ListChapters:
		return c.svc.ListChapters(request.Name, request.Part)
	case constants.GetChapter:
		return c.svc.GetChapter(request.Name, request.Part, request.Chapter)
	case constants.GetArticle:
		return nil, fmt.Errorf("%w: article %s of chapter %s is not available in book '%s'",
			models.ErrNotFound, request.Article, request.Chapter, request.Name)
	}

	return nil, fmt.Errorf("%w: unsupported request type %d", models.ErrBadRequest, request.Type)
}
