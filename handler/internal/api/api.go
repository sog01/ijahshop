package internal

import (
	"github.com/sog01/ijahshop/module"
)

// API is main entity
type API struct {
	mod module.Module
}

// New to create new instance of API main entity
func New(mod module.Module) API {
	return API{
		mod: mod,
	}
}
