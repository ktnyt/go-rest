package rest

import "net/http"

// Handler defines a set of REST endpoint handlers.
type Handler interface {
	// Browse multiple objects filtered by the URL parameters. (GET)
	Browse(http.ResponseWriter, *http.Request)

	// Delete multiple objects filtered by the URL parameters. (DELETE)
	Delete(http.ResponseWriter, *http.Request)

	// Create and store a new object. (POST)
	Create(http.ResponseWriter, *http.Request)

	// Select an object identified by a primary key. (GET)
	Select(http.ResponseWriter, *http.Request)

	// Remove an object identified by a primary key. (DELETE)
	Remove(http.ResponseWriter, *http.Request)

	// Update an entire object identified by a primary key. (PUT)
	Update(http.ResponseWriter, *http.Request)

	// Modify part of an object identified by a primary key. (PATCH)
	Modify(http.ResponseWriter, *http.Request)
}
