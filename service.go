package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ServiceError is an error with an associated HTTP status code.
type ServiceError struct {
	Err  error
	Code int
}

// NewServiceError creates a new ServiceError object.
func NewServiceError(err error, code int) ServiceError {
	return ServiceError{Err: err, Code: code}
}

// Error satisfies the error interface.
func (e ServiceError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Err.Error())
}

// HandleError handles the given error gracefully.
func HandleError(err error, w http.ResponseWriter) bool {
	switch err := err.(type) {
	case ServiceError:
		http.Error(w, err.Error(), err.Code)
		return true
	case nil:
		return false
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
}

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

// PK is the default primary key.
const PK = "pk"

// ServiceBuilder will construct a Service given no arguments.
type ServiceBuilder func() Service

type serviceInterface struct {
	service Service
	pkparam string
}

// NewServiceInterface creates an Interface wrapped around the given Service.
func NewServiceInterface(service Service) Interface {
	return NewServiceInterfaceWithPKParam(service, PK)
}

// NewServiceInterfaceWithPKParam creates a ServiceInterface for the given
// primary key parameter value.
func NewServiceInterfaceWithPKParam(service Service, pkparam string) Interface {
	return serviceInterface{service: service, pkparam: pkparam}
}

func (i serviceInterface) Browse(w http.ResponseWriter, r *http.Request) {
	ctx := InjectParams(r.Context(), r.URL.Query())
	list, err := i.service.Browse(ctx)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(list), w)
}

func (i serviceInterface) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := InjectParams(r.Context(), r.URL.Query())
	list, err := i.service.Delete(ctx)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(list), w)
}

func (i serviceInterface) Create(w http.ResponseWriter, r *http.Request) {
	item, err := i.service.Create(r.Body)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(item), w)
}

func (i serviceInterface) Select(w http.ResponseWriter, r *http.Request) {
	pk := r.Context().Value(i.pkparam).(string)
	item, err := i.service.Select(pk)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(item), w)
}

func (i serviceInterface) Remove(w http.ResponseWriter, r *http.Request) {
	pk := r.Context().Value(i.pkparam).(string)
	item, err := i.service.Remove(pk)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(item), w)
}

func (i serviceInterface) Update(w http.ResponseWriter, r *http.Request) {
	pk := r.Context().Value(i.pkparam).(string)
	item, err := i.service.Update(pk, r.Body)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(item), w)
}

func (i serviceInterface) Modify(w http.ResponseWriter, r *http.Request) {
	pk := r.Context().Value(i.pkparam).(string)
	item, err := i.service.Modify(pk, r.Body)
	if HandleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	HandleError(json.NewEncoder(w).Encode(item), w)
}
