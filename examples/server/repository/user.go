package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

var _ Repository = (*UserRepository)(nil)

type UserRepository struct{}

func (r UserRepository) GetAll() (models []model.Model) {
	for _, m := range database.Storage[database.UserNamespace] {
		models = append(models, m)
	}
	return
}

func (r UserRepository) GetById(id uint) (model model.Model, found bool) {
	model, found = database.Storage[database.UserNamespace][id]
	return
}

func (r UserRepository) Create(model model.Model) model.Model {
	newID := uint(len(database.Storage[database.UserNamespace]))
	model.SetID(newID)
	database.Storage[database.UserNamespace][newID] = model
	return model
}

func (r UserRepository) Delete(id uint) {
	delete(database.Storage[database.UserNamespace], id)
}

func User() Repository {
	return UserRepository{}
}
