package rest

import (
	"context"
	"net/url"
)

// Param is a type alias for specifying URL parameters in contexts.
type Param string

// Params is the key for the URL parameters.
const Params = Param("params")

// InjectParams injects the URL parameters into the given context.
func InjectParams(ctx context.Context, params url.Values) context.Context {
	return context.WithValue(ctx, Params, params)
}
