package rest

// Model defines an interface for a storable model.
type Model interface {
	// Validate the model and return an error if the model is invalid.
	Validate() error

	// MakeKey creates a new key given an integer.
	MakeKey(int) string

	// Merge another interface into this Model.
	Merge(interface{}) error
}

// ModelBuilder will construct a Model given no arguments.
type ModelBuilder func() Model
