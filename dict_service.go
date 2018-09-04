package rest

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func init() {
	gob.Register(&DictService{})
}

// KeyError represents an existing or missing key error.
type KeyError struct {
	key     string
	missing bool
}

// NewKeyError creates a new KeyError.
func NewKeyError(key string, missing bool) error {
	return KeyError{key: key, missing: missing}
}

// Error satisfies the error interface.
func (e KeyError) Error() string {
	if e.missing {
		return fmt.Sprintf("key '%s' does not exist", e.key)
	}
	return fmt.Sprintf("key '%s' already exists", e.key)
}

// FilterFactory creates a filter from the given url parameters.
type FilterFactory func(context.Context) Filter

// Converter takes and interface and convets it to a Model.
type Converter func(interface{}) Model

// DictService provides a Dict service interface.
type DictService struct {
	Dict  *Dict
	Count int

	build   ModelBuilder
	factory FilterFactory
	convert Converter
}

// NewDictService returns a new Dict service.
func NewDictService(build ModelBuilder, factory FilterFactory, convert Converter) Service {
	return &DictService{
		Dict:    NewDict(),
		Count:   0,
		build:   build,
		factory: factory,
		convert: convert,
	}
}

// Browse Dict values filtered by URL parameters.
func (s *DictService) Browse(ctx context.Context) ([]Model, error) {
	indices := s.Dict.Search(s.factory(ctx))
	list := make([]Model, len(indices))
	for j, i := range indices {
		list[j] = s.convert(s.Dict.Values[i])
	}
	return list, nil
}

// Delete Dict values filtered by URL parameters.
func (s *DictService) Delete(ctx context.Context) ([]Model, error) {
	indices := s.Dict.Search(s.factory(ctx))

	keys := make([]string, len(indices))
	for j, i := range indices {
		keys[j] = s.Dict.Keys[i]
	}

	list := make([]Model, len(indices))

	for i, key := range keys {
		if value := s.Dict.Remove(key); value != nil {
			list[i] = s.convert(value)
		}
	}

	return list, nil
}

// Create and store a new value.
func (s *DictService) Create(reader io.Reader) (Model, error) {
	model := s.build()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	key := model.MakeKey(s.Count)
	if err := model.Validate(); err != nil {
		return nil, err
	}

	if !s.Dict.Insert(key, model) {
		err := NewKeyError(key, false)
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	s.Count++

	return model, nil
}

// Select a value identified by the given key.
func (s *DictService) Select(key string) (Model, error) {
	value := s.Dict.Get(key)
	if value == nil {
		err := NewKeyError(key, false)
		return nil, NewServiceError(err, http.StatusBadRequest)
	}
	return s.convert(value), nil
}

// Remove a value identified by the given key.
func (s *DictService) Remove(key string) (Model, error) {
	value := s.Dict.Remove(key)
	if value == nil {
		err := NewKeyError(key, false)
		return nil, NewServiceError(err, http.StatusBadRequest)
	}
	return s.convert(value), nil
}

// Update an entire value identified by the given key.
func (s *DictService) Update(key string, reader io.Reader) (Model, error) {
	value := s.build()
	if err := json.NewDecoder(reader).Decode(&value); err != nil {
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	if err := value.Validate(); err != nil {
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	if !s.Dict.Set(key, value) {
		err := NewKeyError(key, false)
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	return value, nil
}

// Modify part of a value identified by the given key.
func (s *DictService) Modify(key string, reader io.Reader) (Model, error) {
	model := s.build()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	index := s.Dict.Index(key)
	if index == s.Dict.Len() || s.Dict.Keys[index] != key {
		err := NewKeyError(key, false)
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	value := s.convert(s.Dict.Values[index])

	if err := value.Merge(model); err != nil {
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	if err := value.Validate(); err != nil {
		return nil, NewServiceError(err, http.StatusBadRequest)
	}

	s.Dict.Values[index] = value

	return value, nil
}
