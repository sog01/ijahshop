package internal

import (
	"context"
	"database/sql"
)

// ProductAvgValue is entity of product with average value
type ProductAvgValue struct {
	Product
	AverageCost int `db:"average_cost" json:"average_cost"`
	Total       int `json:"total"`
}

// GetProductAvgValue is to get report of product with average value
func (intr Internal) GetProductAvgValue(ctx context.Context) ([]ProductAvgValue, error) {
	var (
		productsWithAvgValue []ProductAvgValue
		query                string
	)

	query = `
	SELECT 
		product.product_id as product_id,
		product.sku as sku,
		product.name as name,
		product.stock as stock,				
		COALESCE(ROUND(AVG(purchase.cost)), 0) as average_cost		
	FROM product
	LEFT JOIN purchase ON product.product_id = purchase.product_id	
	GROUP BY product.product_id
	`

	db := intr.Storage.DB
	err := db.SelectContext(ctx, &productsWithAvgValue, query)
	if err != nil {
		return nil, err
	}

	return productsWithAvgValue, nil

}

// GetProductAvgValueByProductID is to get report of product with average value by productID
func (intr Internal) GetProductAvgValueByProductID(ctx context.Context, productID int64) (ProductAvgValue, error) {
	var (
		productWithAvgValue ProductAvgValue
		query               string
	)

	query = `
	SELECT 
		product.product_id as product_id,
		product.sku as sku,
		product.name as name,
		product.stock as stock,				
		COALESCE(ROUND(AVG(purchase.cost)), 0) as average_cost		
	FROM product	
	LEFT JOIN purchase ON product.product_id = purchase.product_id	
	WHERE 
		product.product_id = ?
	GROUP BY product.product_id
	`

	db := intr.Storage.DB
	row := db.QueryRowxContext(ctx, db.Rebind(query), productID)
	err := row.StructScan(&productWithAvgValue)

	// pass if sql no rows error
	if err == sql.ErrNoRows {
		err = nil
	}

	return productWithAvgValue, err

}
