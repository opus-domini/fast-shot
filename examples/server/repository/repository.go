package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/model"
)

type Repository interface {
	GetAll() (models []model.Model)
	GetById(id uint) (model model.Model, found bool)
	Create(model model.Model) (created model.Model)
	Delete(id uint)
}
