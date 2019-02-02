package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/sog01/ijahshop/storage"
)

// Iinternal is internal contract
type Iinternal interface {
	GetClientByID(ctx context.Context, ID int64) (Client, error)
	CreateClient(ctx context.Context, data Client) error

	// Product function
	GetProduct(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, ID int64) (Product, error)
	StoreProduct(ctx context.Context, tx *sql.Tx, product Product) (ID int64, err error)
	DeleteProduct(ctx context.Context, ID int64) error

	// Purchase Function
	GetPurchaseWithProduct(ctx context.Context) ([]PurchaseWithProduct, error)
	GetPurchaseWithProductByDate(ctx context.Context, dateStart, dateEnd time.Time) ([]PurchaseWithProduct, error)
	GetPurchaseWithProductByID(ctx context.Context, ID int64) (PurchaseWithProduct, error)
	StorePurchase(ctx context.Context, tx *sql.Tx, purchase Purchase) (ID int64, err error)
	StorePurchaseDtl(ctx context.Context, tx *sql.Tx, purchaseDtl PurchaseDtl) (ID int64, err error)

	// Order Function
	GetOrderWithProduct(ctx context.Context) ([]OrderWithProduct, error)
	GetOrderWithProductByDate(ctx context.Context, dateStart, dateEnd time.Time) ([]OrderWithProduct, error)
	GetOrderWithProductByID(ctx context.Context, ID int64) (OrderWithProduct, error)
	StoreOrder(ctx context.Context, tx *sql.Tx, order Order) (ID int64, err error)

	// Report function
	GetProductAvgValue(ctx context.Context) ([]ProductAvgValue, error)
	GetProductAvgValueByProductID(ctx context.Context, productID int64) (ProductAvgValue, error)
}

// Internal is entity of package internal
// Internal is protected package that only can be imported by module package
type Internal struct {
	Storage storage.Storage
}

// New to create instance of internal package
func New(storage storage.Storage) Iinternal {
	return Internal{
		Storage: storage,
	}
}
