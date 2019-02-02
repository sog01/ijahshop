package module

import (
	"github.com/sog01/ijahshop/module/internal"
	"github.com/sog01/ijahshop/storage"
)

// Module is entity of package module
// Module used to define a functionality of services
// commonly used, to be exported into handler
type Module struct {
	Storage  storage.Storage
	internal internal.Iinternal
}

// New to create new instance of module
func New(storage storage.Storage) Module {
	internal := internal.New(storage)
	return Module{
		Storage:  storage,
		internal: internal,
	}
}

// ImportExcelToDB is to import data from excel to database
func (mod Module) ImportExcelToDB(filename string) error {
	return mod.Storage.SeedProductFromEXCEL("files/" + filename)
}
