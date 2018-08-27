package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

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
type FilterFactory func(url.Values) Filter

type dictService struct {
	dict      *Dict
	count     int
	construct Constructor
	factory   FilterFactory
}

// NewDictService returns a new Dict service.
func NewDictService(construct Constructor, factory FilterFactory) Service {
	return &dictService{
		dict:  NewDict(),
		count: 0,

		construct: construct,
		factory:   factory,
	}
}

func (s *dictService) Browse(params url.Values) ([]Model, error) {
	indices := s.dict.Search(s.factory(params))
	list := make([]Model, len(indices))
	for j, i := range indices {
		list[j] = s.dict.Values[i]
	}
	return list, nil
}

func (s *dictService) Delete(params url.Values) ([]Model, error) {
	indices := s.dict.Search(s.factory(params))

	if len(indices) == s.dict.Length() {
		list := s.dict.Values
		s.dict.Clear()
		return list, nil
	}

	keys := make([]string, len(indices))
	for j, i := range indices {
		keys[j] = s.dict.Keys[i]
	}

	list := make([]Model, len(indices))

	for i, key := range keys {
		model := s.dict.Remove(key)
		if model == nil {
			return nil, NewKeyError(key, true)
		}
		list[i] = model
	}

	return list, nil
}

func (s *dictService) Create(reader io.Reader) (Model, error) {
	model := s.construct()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}

	key := model.MakeKey(s.count)
	if err := model.Validate(); err != nil {
		return nil, err
	}

	if !s.dict.Insert(key, model) {
		return nil, NewKeyError(key, false)
	}

	s.count++

	return model, nil
}

func (s *dictService) Select(key string) (Model, error) {
	model := s.dict.Get(key)
	if model == nil {
		return nil, NewKeyError(key, true)
	}
	return model, nil
}

func (s *dictService) Remove(key string) (Model, error) {
	model := s.dict.Remove(key)
	if model == nil {
		return nil, NewKeyError(key, true)
	}
	return model, nil
}

func (s *dictService) Update(key string, reader io.Reader) (Model, error) {
	model := s.construct()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}

	if err := model.Validate(); err != nil {
		return nil, err
	}

	if !s.dict.Set(key, model) {
		return nil, NewKeyError(key, true)
	}

	return model, nil
}

func (s *dictService) Modify(key string, reader io.Reader) (Model, error) {
	model := s.construct()
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}

	index := s.dict.Index(key)
	if index == s.dict.Length() || s.dict.Keys[index] != key {
		return nil, fmt.Errorf("key '%s' does not exist", key)
	}

	value := s.dict.Values[index]

	if err := value.Merge(model); err != nil {
		return nil, err
	}

	if err := value.Validate(); err != nil {
		return nil, err
	}

	s.dict.Values[index] = value

	return value, nil
}
