package model

import (
	"fmt"
)

var _ Model = (*Tuple)(nil)

type (
	// Tuple represents a key-value pair
	Tuple struct {
		ID    uint   `json:"id"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

func (t *Tuple) SetID(id uint) {
	t.ID = id
}

func (t *Tuple) GetID() uint {
	return t.ID
}

func (t *Tuple) String() string {
	return fmt.Sprintf(
		"Tuple{ID: %d, Key: %s, Value: %s}",
		t.ID,
		t.Key,
		t.Value,
	)
}
