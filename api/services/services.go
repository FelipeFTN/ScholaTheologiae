package services

import "scholatheologiae-api/data"

type Services struct {
	data *data.Data
}

func New(data *data.Data) *Services {
	return &Services{
		data: data,
	}
}
