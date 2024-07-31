package repository

import (
	"github.com/opus-domini/fast-shot/examples/server/database"
	"github.com/opus-domini/fast-shot/examples/server/model"
)

type (
	// Repository is the interface that defines the methods to interact with the database.
	Repository interface {
		GetAll() (models []model.Model)
		GetById(id uint) (model model.Model, found bool)
		Create(model model.Model) (created model.Model)
		Delete(id uint)
	}

	// Provider is a struct that contains all the repositories.
	Provider struct {
		User     Repository
		Resource Repository
		Tuple    Repository
	}
)

// NewProvider creates a new provider with the given state.
func NewProvider(state *database.State) *Provider {
	userNamespace := (*state)[database.UserNamespace]
	resourceNamespace := (*state)[database.ResourceNamespace]
	tupleNamespace := (*state)[database.TupleNamespace]

	return &Provider{
		User:     newRepositoryImplementation(&userNamespace),
		Resource: newRepositoryImplementation(&resourceNamespace),
		Tuple:    newRepositoryImplementation(&tupleNamespace),
	}
}
