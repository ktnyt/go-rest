package rest

import (
	"context"
	"io"

	"github.com/pkg/errors"
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
	return ioService{handler: handler}
}

func (s ioService) Browse(ctx context.Context) ([]Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Browse")
	}

	return service.Browse(ctx)
}

func (s ioService) Delete(ctx context.Context) ([]Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Delete")
	}

	list, err := service.Delete(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Delete")
	}

	if err := s.handler.Save(service); err != nil {
		return nil, errors.Wrap(err, "in IO Service Delete")
	}

	return list, nil
}

func (s ioService) Create(reader io.Reader) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Create")
	}

	model, err := service.Create(reader)
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Create")
	}

	if err := s.handler.Save(service); err != nil {
		return nil, errors.Wrap(err, "in IO Service Create")
	}

	return model, nil
}

func (s ioService) Select(key string) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Select")
	}

	return service.Select(key)
}

func (s ioService) Remove(key string) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Remove")
	}

	model, err := service.Remove(key)
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Remove")
	}

	if err := s.handler.Save(service); err != nil {
		return nil, errors.Wrap(err, "in IO Service Remove")
	}

	return model, nil
}

func (s ioService) Update(key string, reader io.Reader) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Update")
	}

	model, err := service.Update(key, reader)
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Update")
	}

	if err := s.handler.Save(service); err != nil {
		return nil, errors.Wrap(err, "in IO Service Update")
	}

	return model, nil
}

func (s ioService) Modify(key string, reader io.Reader) (Model, error) {
	service, err := s.handler.Load()
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Modify")
	}

	model, err := service.Modify(key, reader)
	if err != nil {
		return nil, errors.Wrap(err, "in IO Service Modify")
	}

	if err := s.handler.Save(service); err != nil {
		return nil, errors.Wrap(err, "in IO Service Modify")
	}

	return model, nil
}
