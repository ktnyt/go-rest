package rest

import (
	"io"
	"net/url"
)

// Service defines interfaces for manipulating values for a persistence backend.
type Service interface {
	Browse(url.Values) ([]Model, error)
	Delete(url.Values) ([]Model, error)
	Create(io.Reader) (Model, error)
	Select(string) (Model, error)
	Remove(string) (Model, error)
	Update(string, io.Reader) (Model, error)
	Modify(string, io.Reader) (Model, error)
}
