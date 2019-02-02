package internal

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// Purchase is entity that represent schema on table purchase
type Purchase struct {
	PurchaseID       int64     `db:"purchase_id" json:"purchase_id"`
	ProductID        int64     `db:"product_id" json:"product_id"`
	QuantityOrder    int       `db:"quantity_order" json:"quantity_order"`
	QuantityAccepted int       `db:"quantity_accepted" json:"quantity_accepted"`
	Description      string    `db:"description" json:"description"`
	InvoiceNumber    string    `db:"invoice_number" json:"invoice_number"`
	Cost             int64     `db:"cost" json:"cost"`
	Date             time.Time `db:"date" json:"-"`
	DateStr          string    `db:"date_str" json:"date"`
	IsFinish         bool      `db:"is_finish" json:"is_finish"`
	Total            int64     `json:"total"`
}

// PurchaseDtl is entity that represent schema on table purchase_detail
type PurchaseDtl struct {
	PurchaseDtlID int64     `db:"purchase_detail_id" json:"purchase_detail_id"`
	PurchaseID    int64     `db:"purchase_id" json:"purchase_id"`
	Quantity      int       `db:"quantity" json:"quantity"`
	Description   string    `db:"description" json:"description"`
	Date          time.Time `db:"date" json:"date"`
}

// PurchaseWithProduct is entity of purchase with product
type PurchaseWithProduct struct {
	Purchase
	Product Product `json:"product"`
}

// this is a main query. it will be used on many place
// so to reduce redudancy, this query need to be declared as a global variable
var qSelectPurchase = `
	SELECT 
		purchase.purchase_id,
		purchase.product_id,
		purchase.quantity_order,
		purchase.quantity_accepted,
		purchase.description,
		purchase.invoice_number,
		purchase.cost,
		purchase.date,
		purchase.is_finish,
		product.name,
		product.sku,
		product.stock
	FROM purchase
	JOIN product ON purchase.product_id = product.product_id
`

// GetPurchaseWithProduct is used to get purchased with product
func (intr Internal) GetPurchaseWithProduct(ctx context.Context) ([]PurchaseWithProduct, error) {
	var (
		purchaseWithProducts []PurchaseWithProduct
		query                string
	)

	query = qSelectPurchase

	db := intr.Storage.DB
	row, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		purchaseWithProduct := PurchaseWithProduct{}
		err := row.Scan(
			&purchaseWithProduct.PurchaseID,
			&purchaseWithProduct.ProductID,
			&purchaseWithProduct.QuantityOrder,
			&purchaseWithProduct.QuantityAccepted,
			&purchaseWithProduct.Description,
			&purchaseWithProduct.InvoiceNumber,
			&purchaseWithProduct.Cost,
			&purchaseWithProduct.DateStr,
			&purchaseWithProduct.IsFinish,
			&purchaseWithProduct.Product.Name,
			&purchaseWithProduct.Product.Sku,
			&purchaseWithProduct.Product.Stock,
		)

		if err != nil {
			return nil, err
		}

		purchaseWithProduct.Product.ProductID = purchaseWithProduct.ProductID

		// convert date string into date time.Time
		// spit date string to remove character +00:00
		splitDateStr := strings.Split(purchaseWithProduct.DateStr, "+")
		DateStr := strings.Trim(splitDateStr[0], " ")
		purchaseWithProduct.Date, err = time.Parse("2006-01-02 15:04:05", DateStr)
		if err != nil {
			return nil, err
		}

		// using date format: yyyy-MM-dd HH:mm:ss
		// to standarize date convenient
		purchaseWithProduct.DateStr = purchaseWithProduct.Date.Format("2006-01-02 15:04:05")

		purchaseWithProducts = append(purchaseWithProducts, purchaseWithProduct)

	}

	return purchaseWithProducts, nil
}

// GetPurchaseWithProductByDate is used to get purchased with product by date
func (intr Internal) GetPurchaseWithProductByDate(ctx context.Context, dateStart, dateEnd time.Time) ([]PurchaseWithProduct, error) {
	var (
		purchaseWithProducts []PurchaseWithProduct
		query                string
	)

	dateStartMidnight := time.Date(
		dateStart.Year(),
		time.Month(dateStart.Month()),
		dateStart.Day(),
		0, 0, 0, 0, time.UTC,
	)

	dateEndMidnight := time.Date(
		dateEnd.Year(),
		time.Month(dateEnd.Month()),
		dateEnd.Day(),
		23, 59, 59, 0, time.UTC,
	)

	query = qSelectPurchase
	query += `WHERE
		purchase.date > ? AND purchase.date <= ?
	`

	db := intr.Storage.DB
	row, err := db.QueryxContext(ctx, query, dateStartMidnight, dateEndMidnight)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		purchaseWithProduct := PurchaseWithProduct{}
		err := row.Scan(
			&purchaseWithProduct.PurchaseID,
			&purchaseWithProduct.ProductID,
			&purchaseWithProduct.QuantityOrder,
			&purchaseWithProduct.QuantityAccepted,
			&purchaseWithProduct.Description,
			&purchaseWithProduct.InvoiceNumber,
			&purchaseWithProduct.Cost,
			&purchaseWithProduct.DateStr,
			&purchaseWithProduct.IsFinish,
			&purchaseWithProduct.Product.Name,
			&purchaseWithProduct.Product.Sku,
			&purchaseWithProduct.Product.Stock,
		)

		if err != nil {
			return nil, err
		}

		purchaseWithProduct.Product.ProductID = purchaseWithProduct.ProductID

		// convert date string into date time.Time
		// spit date string to remove character +00:00
		// date format: yyyy-MM-dd HH:mm:ss
		splitDateStr := strings.Split(purchaseWithProduct.DateStr, "+")
		DateStr := strings.Trim(splitDateStr[0], " ")
		purchaseWithProduct.Date, err = time.Parse("2006-01-02 15:04:05", DateStr)
		if err != nil {
			return nil, err
		}

		// using date format: yyyy-MM-dd HH:mm:ss
		// to standarize date convenient
		purchaseWithProduct.DateStr = purchaseWithProduct.Date.Format("2006-01-02 15:04:05")

		purchaseWithProducts = append(purchaseWithProducts, purchaseWithProduct)

	}

	return purchaseWithProducts, nil
}

// GetPurchaseWithProductByID is used to get purchased with product by ID
func (intr Internal) GetPurchaseWithProductByID(ctx context.Context, ID int64) (PurchaseWithProduct, error) {
	var (
		purchaseWithProduct PurchaseWithProduct
		query               string
	)
	query = qSelectPurchase
	query += `WHERE
				purchase.purchase_id = ?
			`

	db := intr.Storage.DB
	row := db.QueryRowxContext(ctx, query, ID)
	err := row.Scan(
		&purchaseWithProduct.PurchaseID,
		&purchaseWithProduct.ProductID,
		&purchaseWithProduct.QuantityOrder,
		&purchaseWithProduct.QuantityAccepted,
		&purchaseWithProduct.Description,
		&purchaseWithProduct.InvoiceNumber,
		&purchaseWithProduct.Cost,
		&purchaseWithProduct.DateStr,
		&purchaseWithProduct.IsFinish,
		&purchaseWithProduct.Product.Name,
		&purchaseWithProduct.Product.Sku,
		&purchaseWithProduct.Product.Stock,
	)

	// keep returning value but with empty struct
	// since no rows is not error in a system
	if err == sql.ErrNoRows {
		return PurchaseWithProduct{}, nil
	}

	if err != nil {
		return PurchaseWithProduct{}, err
	}

	purchaseWithProduct.Product.ProductID = purchaseWithProduct.ProductID

	// convert date string into date time.Time
	// spit date string to remove character +00:00
	// date format: yyyy-MM-dd HH:mm:ss
	splitDateStr := strings.Split(purchaseWithProduct.DateStr, "+")
	DateStr := strings.Trim(splitDateStr[0], " ")
	purchaseWithProduct.Date, err = time.Parse("2006-01-02 15:04:05", DateStr)
	if err != nil {
		return PurchaseWithProduct{}, err
	}

	// using date format: yyyy-MM-dd HH:mm:ss
	// to standarize date convenient
	purchaseWithProduct.DateStr = purchaseWithProduct.Date.Format("2006-01-02 15:04:05")

	return purchaseWithProduct, nil
}

// StorePurchase is to store purchase into database
func (intr Internal) StorePurchase(ctx context.Context, tx *sql.Tx, purchase Purchase) (ID int64, err error) {
	var args []interface{}
	query := `INSERT INTO purchase  
					(
						product_id,
						quantity_order,
						quantity_accepted,
						description,
						invoice_number,
						cost,
						date,
						is_finish 
					)
			VALUES (
						?, 
						?, 
						?,
						?,
						?,
						?,
						?,
						?						
					)
			`
	args = append(args,
		purchase.ProductID,
		purchase.QuantityOrder,
		purchase.QuantityAccepted,
		purchase.Description,
		purchase.InvoiceNumber,
		purchase.Cost,
		purchase.Date,
		purchase.IsFinish,
	)

	if purchase.PurchaseID != 0 {
		query = `UPDATE purchase 
				 SET 
						product_id = ?,
						quantity_order = ?,
						quantity_accepted = ?,
						description = ?,
						invoice_number = ?,
						cost = ?,
						date = ?,
						is_finish = ?
				WHERE 
						purchase_id = ?
						 
		`
		args = append(args, purchase.PurchaseID)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	if purchase.PurchaseID == 0 {
		// no need to check error, since it will be occurred by database incompatibility
		purchase.PurchaseID, _ = result.LastInsertId()
	}

	return purchase.PurchaseID, err
}

// StorePurchaseDtl is to store purchase detail into database
func (intr Internal) StorePurchaseDtl(ctx context.Context, tx *sql.Tx, purchaseDtl PurchaseDtl) (ID int64, err error) {
	var args []interface{}
	query := `INSERT INTO purchase_detail  
					(
						purchase_id,
						quantity,
						description,												
						date						
					)
			VALUES (
						?, 
						?, 
						?,
						?												
					)
			`
	args = append(args,
		purchaseDtl.PurchaseID,
		purchaseDtl.Quantity,
		purchaseDtl.Description,
		purchaseDtl.Date,
	)
	if purchaseDtl.PurchaseDtlID != 0 {
		query = `UPDATE purchase_detail 
				 SET 
						purchase_id = ?,
						quantity = ?,
						description = ?,												
						date = ?
				WHERE 
						purchase_detail_id = ?
						 
		`
		args = append(args, purchaseDtl.PurchaseDtlID)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	if purchaseDtl.PurchaseID == 0 {
		// no need to check error, since it will be occurred by database incompatibility
		purchaseDtl.PurchaseID, _ = result.LastInsertId()
	}

	return purchaseDtl.PurchaseID, err
}
