package rest

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
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

// DictService provides a Dict service interface.
type DictService struct {
	Dict  *Dict
	Count int

	build   ModelBuilder
	factory FilterFactory
}

// NewDictService returns a new Dict service.
func NewDictService(build ModelBuilder, factory FilterFactory) Service {
	return &DictService{
		Dict:    NewDict(),
		Count:   0,
		build:   build,
		factory: factory,
	}
}

// Browse Dict values filtered by URL parameters.
func (s *DictService) Browse(ctx context.Context) ([]Model, error) {
	indices := s.Dict.Search(s.factory(ctx))
	list := make([]Model, len(indices))
	for j, i := range indices {
		list[j] = s.Dict.Values[i]
	}
	return list, nil
}

// Delete Dict values filtered by URL parameters.
func (s *DictService) Delete(ctx context.Context) ([]Model, error) {
	indices := s.Dict.Search(s.factory(ctx))

	if len(indices) == s.Dict.Len() {
		list := s.Dict.Values
		s.Dict.Clear()
		return list, nil
	}

	keys := make([]string, len(indices))
	for j, i := range indices {
		keys[j] = s.Dict.Keys[i]
	}

	list := make([]Model, len(indices))

	for i, key := range keys {
		model := s.Dict.Remove(key)
		if model == nil {
			return nil, NewKeyError(key, true)
		}
		list[i] = model
	}

	return list, nil
}

// Create and store a new value.
func (s *DictService) Create(reader io.Reader) (Model, error) {
	model := s.build()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}

	key := model.MakeKey(s.Count)
	if err := model.Validate(); err != nil {
		return nil, err
	}

	if !s.Dict.Insert(key, model) {
		return nil, NewKeyError(key, false)
	}

	s.Count++

	return model, nil
}

// Select a value identified by the given key.
func (s *DictService) Select(key string) (Model, error) {
	model := s.Dict.Get(key)
	if model == nil {
		return nil, NewKeyError(key, true)
	}
	return model, nil
}

// Remove a value identified by the given key.
func (s *DictService) Remove(key string) (Model, error) {
	model := s.Dict.Remove(key)
	if model == nil {
		return nil, NewKeyError(key, true)
	}
	return model, nil
}

// Update an entire value identified by the given key.
func (s *DictService) Update(key string, reader io.Reader) (Model, error) {
	model := s.build()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}

	if err := model.Validate(); err != nil {
		return nil, err
	}

	if !s.Dict.Set(key, model) {
		return nil, NewKeyError(key, true)
	}

	return model, nil
}

// Modify part of a value identified by the given key.
func (s *DictService) Modify(key string, reader io.Reader) (Model, error) {
	model := s.build()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}

	index := s.Dict.Index(key)
	if index == s.Dict.Len() || s.Dict.Keys[index] != key {
		return nil, fmt.Errorf("key '%s' does not exist", key)
	}

	value := s.Dict.Values[index]

	if err := value.Merge(model); err != nil {
		return nil, err
	}

	if err := value.Validate(); err != nil {
		return nil, err
	}

	s.Dict.Values[index] = value

	return value, nil
}
