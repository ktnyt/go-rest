package rest

import (
	"io"
	"net/url"
)

// IOHandler defines interfaces for persisting a service.
type IOHandler interface {
	Save(Service) error
	Load() (Service, error)
}

type ioService struct {
	handler IOHandler
}

// NewIOService returns a new IO service.
func NewIOService(handler IOHandler) Service {
	return &ioService{handler: handler}
}

func (s *ioService) Browse(params url.Values) ([]Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	return service.Browse(params)
}

func (s *ioService) Delete(params url.Values) ([]Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	list, err := service.Delete(params)
	if err != nil {
		return nil, err
	}

	if err := s.handler.Save(service); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *ioService) Create(reader io.Reader) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	model, err := service.Create(reader)
	if err != nil {
		return nil, err
	}

	if err := s.handler.Save(service); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ioService) Select(key string) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	return service.Select(key)
}

func (s *ioService) Remove(key string) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	model, err := service.Remove(key)
	if err != nil {
		return nil, err
	}

	if err := s.handler.Save(service); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ioService) Update(key string, reader io.Reader) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	model, err := service.Update(key, reader)
	if err != nil {
		return nil, err
	}

	if err := s.handler.Save(service); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ioService) Modify(key string, reader io.Reader) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, err
	}

	model, err := service.Modify(key, reader)
	if err != nil {
		return nil, err
	}

	if err := s.handler.Save(service); err != nil {
		return nil, err
	}

	return model, nil
}
