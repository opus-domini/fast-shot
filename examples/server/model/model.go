package model

type (
	// Model represents a generic model in the system
	Model interface {
		SetID(id uint)
		GetID() uint
		String() string
	}
)
