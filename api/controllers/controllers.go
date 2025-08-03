package controllers

import (
	"scholatheologiae-api/data"
	"scholatheologiae-api/services"
)

type Controllers struct {
	data *data.Data
	svc  *services.Services
}

func New(d *data.Data) *Controllers {
	var controllers Controllers

	controllers.data = d
	controllers.svc = services.New(d)

	return &controllers
}
