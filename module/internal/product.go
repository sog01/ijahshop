package internal

import (
	"context"
	"database/sql"
)

// Product is entity that represent schema on table product
type Product struct {
	ProductID int64  `db:"product_id" json:"product_id"`
	Name      string `db:"name" json:"product_name"`
	Sku       string `db:"sku" json:"product_sku"`
	Stock     int    `db:"stock" json:"product_stock"`
}

// this is a main query. it will be used on many place
// so to reduce redudancy, this query need to be declared as a global variable
var qSelectProduct = `
			SELECT 
					product_id,
					name,
					sku,
					stock
			FROM product
			`

// GetProduct is used to get all product
func (intr Internal) GetProduct(ctx context.Context) ([]Product, error) {
	var (
		products []Product
		query    string
	)

	query = qSelectProduct
	db := intr.Storage.DB
	err := db.SelectContext(ctx, &products, query)
	return products, err
}

// GetProductByID is used to get productByID
func (intr Internal) GetProductByID(ctx context.Context, ID int64) (Product, error) {
	var (
		product Product
		query   string
	)

	query = qSelectProduct
	query += `WHERE
				product_id = ?
			`
	db := intr.Storage.DB
	row := db.QueryRowxContext(ctx, db.Rebind(query), ID)
	err := row.StructScan(&product)

	// pass if sql no rows error
	if err == sql.ErrNoRows {
		err = nil
	}
	return product, err
}

// StoreProduct is to store product into database
func (intr Internal) StoreProduct(ctx context.Context, tx *sql.Tx, product Product) (ID int64, err error) {
	var args []interface{}
	query := `INSERT INTO product  
					(
						name,
						sku,
						stock 
					)
			VALUES (
						?, 
						?, 
						?						
					)
			`
	args = append(args, product.Name, product.Sku, product.Stock)
	if product.ProductID != 0 {
		query = `UPDATE product 
				 SET 
						name = ?,
						sku = ?,
						stock = ?
				WHERE 
						product_id = ?
						 
		`
		args = append(args, product.ProductID)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	if product.ProductID == 0 {
		// no need to check error, since it will be occurred by database incompatibility
		product.ProductID, _ = result.LastInsertId()
	}

	return product.ProductID, err
}

// DeleteProduct is to delete product from database by single ID
func (intr Internal) DeleteProduct(ctx context.Context, ID int64) error {
	query := `DELETE FROM product 
			  WHERE 
				  product_id = ?
			 `
	db := intr.Storage.DB

	_, err := db.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}

	return err
}
