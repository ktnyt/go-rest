package rest

import (
	"encoding/gob"
	"os"
)

type fileHandler string

func (h fileHandler) Save(s Service) error {
	f, err := os.Create(string(h))
	if err != nil {
		return err
	}

	return gob.NewEncoder(f).Encode(s)
}

func (h fileHandler) Load() (Service, error) {
	f, err := os.Open(string(h))
	if err != nil {
		return nil, err
	}

	var s Service
	if err := gob.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}

// NewFileService creates a file persisted service.
func NewFileService(filename string) Service {
	return NewIOService(fileHandler(filename))
}
