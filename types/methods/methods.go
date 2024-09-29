package methods

import (
	"net/http"
)

const (
	MGet     = http.MethodGet
	MPost    = http.MethodPost
	MPatch   = http.MethodPatch
	MPut     = http.MethodPut
	MDelete  = http.MethodDelete
	MHead    = http.MethodHead
	MTrace   = http.MethodTrace
	MOptions = http.MethodOptions
)

type (
	Get     struct{}
	Post    struct{}
	Patch   struct{}
	Put     struct{}
	Delete  struct{}
	Head    struct{}
	Trace   struct{}
	Options struct{}

	METHOD interface {
		Post | Get | Patch | Put |
			Delete | Head | Trace | Options
	}
)
