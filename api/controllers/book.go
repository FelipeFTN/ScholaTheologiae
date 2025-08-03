package controllers

import (
	"errors"

	"scholatheologiae-api/constants"
	"scholatheologiae-api/models"
)

func (c *Controllers) Read(request models.BookRequest) (any, error) {
	switch request.Type {
	case constants.TYPE_LIST_PARTS:
		return c.svc.ListParts(request.Name)
	case constants.TYPE_LIST_CHAPTERS:
		return c.svc.ListChapters(request.Name, request.Part)
	case constants.TYPE_GET_CHAPTER:
		return c.svc.GetChapter(request.Name, request.Part, request.Chapter)
	}

	return nil, errors.New("Unexpected error in Book Controller: " + request.Type)
}
