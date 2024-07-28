package model

import (
	"fmt"
	"time"
)

var _ Model = (*Resource)(nil)

type (
	// Resource represents a resource in the system
	Resource struct {
		ID        uint      `json:"id"`
		Type      string    `json:"type"`
		Data      string    `json:"data"`
		Timestamp time.Time `json:"timestamp"`
	}
)

func (r *Resource) SetID(id uint) {
	r.ID = id
}

func (r *Resource) GetID() uint {
	return r.ID
}

func (r *Resource) String() string {
	return fmt.Sprintf(
		"Resource{ID: %d, Type: %s, Data: %s, Timestamp: %s}",
		r.ID,
		r.Type,
		r.Data,
		r.Timestamp.Format(time.DateOnly),
	)
}
