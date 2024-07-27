package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

var _ Repository = (*ResourceRepository)(nil)

type ResourceRepository struct{}

func (r ResourceRepository) GetAll() (models []model.Model) {
	for _, m := range database.Storage[database.ResourceNamespace] {
		models = append(models, m)
	}
	return
}

func (r ResourceRepository) GetById(id uint) (model model.Model, found bool) {
	model, found = database.Storage[database.ResourceNamespace][id]
	return
}

func (r ResourceRepository) Create(model model.Model) model.Model {
	newID := uint(len(database.Storage[database.ResourceNamespace]))
	model.SetID(newID)
	database.Storage[database.ResourceNamespace][newID] = model
	return model
}

func (r ResourceRepository) Delete(id uint) {
	delete(database.Storage[database.ResourceNamespace], id)
}

func Resource() Repository {
	return ResourceRepository{}
}
