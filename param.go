package rest

import (
	"context"
	"net/url"
)

// Param is a type alias for specifying URL parameters in contexts.
type Param string

// InjectParams injects the URL parameters into the given context.
func InjectParams(ctx context.Context, params url.Values) context.Context {
	for key, value := range params {
		ctx = context.WithValue(ctx, Param(key), value)
	}
	return ctx
}
