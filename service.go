package rest

import (
	"context"
	"io"
)

// Service defines interfaces for manipulating values for a persistence backend.
type Service interface {
	Browse(context.Context) ([]Model, error)
	Delete(context.Context) ([]Model, error)
	Create(io.Reader) (Model, error)
	Select(string) (Model, error)
	Remove(string) (Model, error)
	Update(string, io.Reader) (Model, error)
	Modify(string, io.Reader) (Model, error)
}

// ServiceBuilder will construct a Service given no arguments.
type ServiceBuilder func() Service
