package controller

import "scholatheologiae-api/handler"

type Controller struct {
	db *handler.SQLiteHandler
}

func New(db *handler.SQLiteHandler) *Controller {
	var Controller Controller

	Controller.db = db

	return &Controller
}
