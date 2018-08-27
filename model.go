package rest

// Model defines an interface for a storable model.
type Model interface {
	// Validate the model and return an error if the model is invalid.
	Validate() error

	// MakeKey creates a new key given an integer.
	MakeKey(int) string

	// Merge another Model into this Model.
	Merge(Model) error
}

// ModelBuilder will construct a Model given no arguments.
type ModelBuilder func() Model
