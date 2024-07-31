package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

var _ Repository = (*Implementation)(nil)

type Implementation struct {
	table *database.Table
}

func (r Implementation) GetAll() (models []model.Model) {
	for _, m := range *r.table {
		models = append(models, m)
	}
	return
}

func (r Implementation) GetById(id uint) (model model.Model, found bool) {
	model, found = (*r.table)[id]
	return
}

func (r Implementation) Create(model model.Model) model.Model {
	newID := uint(len(*r.table))
	model.SetID(newID)
	(*r.table)[newID] = model
	return model
}

func (r Implementation) Delete(id uint) {
	delete(*r.table, id)
}

func newRepositoryImplementation(table *database.Table) *Implementation {
	return &Implementation{
		table: table,
	}
}
