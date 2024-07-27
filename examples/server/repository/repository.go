package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

type (
	Repository interface {
		GetAll() []model.Model
		GetById(id uint) (model.Model, bool)
		Create(model model.Model)
		Delete(id uint)
	}

	UserRepository     struct{}
	ResourceRepository struct{}
	TupleRepository    struct{}
)

func (r UserRepository) GetAll() []model.Model {
	var models []model.Model
	for _, m := range database.Storage[database.UserNamespace] {
		models = append(models, m)
	}
	return models
}

func (r UserRepository) GetById(id uint) (model.Model, bool) {
	user, ok := database.Storage[database.UserNamespace][id]
	return user, ok
}

func (r UserRepository) Create(model model.Model) {
	newID := uint(len(database.Storage[database.UserNamespace]))
	model.SetID(newID)
	database.Storage[database.UserNamespace][newID] = model
}

func (r UserRepository) Delete(id uint) {
	delete(database.Storage[database.UserNamespace], id)
}

func (r ResourceRepository) GetAll() []model.Model {
	var models []model.Model
	for _, m := range database.Storage[database.ResourceNamespace] {
		models = append(models, m)
	}
	return models
}

func (r ResourceRepository) GetById(id uint) (model.Model, bool) {
	resource, ok := database.Storage[database.ResourceNamespace][id]
	return resource, ok
}

func (r ResourceRepository) Create(model model.Model) {
	newID := uint(len(database.Storage[database.ResourceNamespace]))
	model.SetID(newID)
	database.Storage[database.ResourceNamespace][newID] = model
}

func (r ResourceRepository) Delete(id uint) {
	delete(database.Storage[database.ResourceNamespace], id)
}

func (r TupleRepository) GetAll() []model.Model {
	var models []model.Model
	for _, m := range database.Storage[database.TupleNamespace] {
		models = append(models, m)
	}
	return models
}

func (r TupleRepository) GetById(id uint) (model.Model, bool) {
	tuple, ok := database.Storage[database.TupleNamespace][id]
	return tuple, ok
}

func (r TupleRepository) Create(model model.Model) {
	newID := uint(len(database.Storage[database.TupleNamespace]))
	model.SetID(newID)
	database.Storage[database.TupleNamespace][newID] = model
}

func (r TupleRepository) Delete(id uint) {
	delete(database.Storage[database.TupleNamespace], id)
}

func User() Repository {
	return UserRepository{}
}

func Resource() Repository {
	return ResourceRepository{}
}

func Tuple() Repository {
	return TupleRepository{}
}
