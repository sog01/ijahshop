package module

import (
	"context"

	"github.com/sog01/ijahshop/module/internal"
)

// ReqProduct is entity of inputed product
// use to make request that will be stored into database
type ReqProduct struct {
	internal.Product
}

// GetProduct is used to get all product
func (mod Module) GetProduct(ctx context.Context) ([]internal.Product, error) {

	return mod.internal.GetProduct(ctx)
}

// GetProductByID is used to get product by ID
func (mod Module) GetProductByID(ctx context.Context, ID int64) (internal.Product, error) {

	return mod.internal.GetProductByID(ctx, ID)
}

// StoreProduct is to store product into database
func (mod Module) StoreProduct(ctx context.Context, reqProduct ReqProduct) (ID int64, err error) {

	product := internal.Product{
		ProductID: reqProduct.ProductID,
		Name:      reqProduct.Name,
		Sku:       reqProduct.Sku,
		Stock:     reqProduct.Stock,
	}

	db := mod.Storage.DB
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	ID, err = mod.internal.StoreProduct(ctx, tx, product)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return ID, tx.Commit()
}

// DeleteProduct is to delete product from database
func (mod Module) DeleteProduct(ctx context.Context, ID int64) error {
	err := mod.internal.DeleteProduct(ctx, ID)
	if err != nil {
		return err
	}

	return err
}

// WriteProductToCSV to write product entity to CSV
func (mod Module) WriteProductToCSV(ctx context.Context) error {
	product, err := mod.internal.GetProduct(ctx)
	if err != nil {
		return err
	}
	return mod.writeToCSV(ctx, "Catatan Jumlah Barang", product)
}
