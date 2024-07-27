package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

var _ Repository = (*TupleRepository)(nil)

type TupleRepository struct{}

func (r TupleRepository) GetAll() (models []model.Model) {
	for _, m := range database.Storage[database.TupleNamespace] {
		models = append(models, m)
	}
	return
}

func (r TupleRepository) GetById(id uint) (model model.Model, found bool) {
	model, found = database.Storage[database.TupleNamespace][id]
	return
}

func (r TupleRepository) Create(model model.Model) model.Model {
	newID := uint(len(database.Storage[database.TupleNamespace]))
	model.SetID(newID)
	database.Storage[database.TupleNamespace][newID] = model
	return model
}

func (r TupleRepository) Delete(id uint) {
	delete(database.Storage[database.TupleNamespace], id)
}

func Tuple() Repository {
	return TupleRepository{}
}
