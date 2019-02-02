package internal

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// Order is entity that represent schema on table order
type Order struct {
	OrderID       int64     `db:"order_id" json:"order_id"`
	OrderIDFormat string    `db:"order_id_format" json:"order_id_format"`
	ProductID     int64     `db:"product_id" json:"product_id"`
	Quantity      int       `db:"quantity" json:"quantity"`
	Description   string    `db:"description" json:"description"`
	Date          time.Time `db:"date" json:"-"`
	DateStr       string    `db:"date_str" json:"date"`
	Price         int64     `db:"price" json:"price"`
	Total         int64     `json:"total"`
}

// OrderWithProduct is entity that represent schema on table order
// with product
type OrderWithProduct struct {
	Order
	Product Product `json:"product"`
}

// this is a main query. it will be used on many place
// so, to reduce redudancy, this query need to be declared as a global variable
var qSelectOrder = `
	SELECT 
		orders.order_id,
		orders.order_id_format,
		orders.product_id,
		orders.quantity,
		orders.description,
		orders.date as date_str,
		orders.price,
		product.name,
		product.sku,
		product.stock
	FROM orders
	JOIN product ON orders.product_id = product.product_id
`

// GetOrderWithProduct is used to get all order with product
func (intr Internal) GetOrderWithProduct(ctx context.Context) ([]OrderWithProduct, error) {
	var (
		ordersWithProduct []OrderWithProduct
		query             string
	)

	query = qSelectOrder
	db := intr.Storage.DB
	row, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		orderWithProduct := OrderWithProduct{}
		err = row.Scan(
			&orderWithProduct.OrderID,
			&orderWithProduct.OrderIDFormat,
			&orderWithProduct.ProductID,
			&orderWithProduct.Quantity,
			&orderWithProduct.Description,
			&orderWithProduct.DateStr,
			&orderWithProduct.Price,
			&orderWithProduct.Product.Name,
			&orderWithProduct.Product.Sku,
			&orderWithProduct.Product.Stock,
		)
		if err != nil {
			return nil, err
		}

		orderWithProduct.Product.ProductID = orderWithProduct.ProductID

		// convert date string into date time.Time
		// spit date string to remove character +00:00
		splitDateStr := strings.Split(orderWithProduct.DateStr, "+")
		DateStr := strings.Trim(splitDateStr[0], " ")
		orderWithProduct.Date, err = time.Parse("2006-01-02 15:04:05", DateStr)
		if err != nil {
			return nil, err
		}

		// using date format: yyyy-MM-dd HH:mm:ss
		// to standarize date convenient
		orderWithProduct.DateStr = orderWithProduct.Date.Format("2006-01-02 15:04:05")

		ordersWithProduct = append(ordersWithProduct, orderWithProduct)

	}

	return ordersWithProduct, err
}

// GetOrderWithProductByDate is used to get all order with product by filter date
func (intr Internal) GetOrderWithProductByDate(ctx context.Context, dateStart, dateEnd time.Time) ([]OrderWithProduct, error) {
	var (
		ordersWithProduct []OrderWithProduct
		query             string
	)

	query = qSelectOrder
	query += `WHERE
		orders.date > ? AND orders.date <= ?
	`

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

	db := intr.Storage.DB
	row, err := db.QueryxContext(ctx, query, dateStartMidnight, dateEndMidnight)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		orderWithProduct := OrderWithProduct{}
		err = row.Scan(
			&orderWithProduct.OrderID,
			&orderWithProduct.OrderIDFormat,
			&orderWithProduct.ProductID,
			&orderWithProduct.Quantity,
			&orderWithProduct.Description,
			&orderWithProduct.DateStr,
			&orderWithProduct.Price,
			&orderWithProduct.Product.Name,
			&orderWithProduct.Product.Sku,
			&orderWithProduct.Product.Stock,
		)
		if err != nil {
			return nil, err
		}

		orderWithProduct.Product.ProductID = orderWithProduct.ProductID

		// convert date string into date time.Time
		// spit date string to remove character +00:00
		// date format: yyyy-MM-dd HH:mm:ss
		splitDateStr := strings.Split(orderWithProduct.DateStr, "+")
		DateStr := strings.Trim(splitDateStr[0], " ")
		orderWithProduct.Date, err = time.Parse("2006-01-02 15:04:05", DateStr)
		if err != nil {
			return nil, err
		}

		// using date format: yyyy-MM-dd HH:mm:ss
		// to standarize date convenient
		orderWithProduct.DateStr = orderWithProduct.Date.Format("2006-01-02 15:04:05")

		ordersWithProduct = append(ordersWithProduct, orderWithProduct)

	}

	return ordersWithProduct, err
}

// GetOrderWithProductByID is used to get order with product by ID
func (intr Internal) GetOrderWithProductByID(ctx context.Context, ID int64) (OrderWithProduct, error) {
	var (
		orderWithProduct OrderWithProduct
		query            string
	)
	query = qSelectOrder
	query += `WHERE
				order_id = ?
			`
	db := intr.Storage.DB
	row := db.QueryRowxContext(ctx, query, ID)
	err := row.Scan(
		&orderWithProduct.OrderID,
		&orderWithProduct.OrderIDFormat,
		&orderWithProduct.ProductID,
		&orderWithProduct.Quantity,
		&orderWithProduct.Description,
		&orderWithProduct.DateStr,
		&orderWithProduct.Price,
		&orderWithProduct.Product.Name,
		&orderWithProduct.Product.Sku,
		&orderWithProduct.Product.Stock,
	)

	// keep returning value but with empty struct
	// since no rows is not error in a system
	if err == sql.ErrNoRows {
		return OrderWithProduct{}, nil
	}

	if err != nil {
		return OrderWithProduct{}, err
	}

	orderWithProduct.Product.ProductID = orderWithProduct.ProductID

	// convert date string into date time.Time
	// spit date string to remove character +00:00
	// date format: yyyy-MM-dd HH:mm:ss
	splitDateStr := strings.Split(orderWithProduct.DateStr, "+")
	DateStr := strings.Trim(splitDateStr[0], " ")
	orderWithProduct.Date, err = time.Parse("2006-01-02 15:04:05", DateStr)
	if err != nil {
		return OrderWithProduct{}, err
	}

	// using date format: yyyy-MM-dd HH:mm:ss
	// to standarize date convenient
	orderWithProduct.DateStr = orderWithProduct.Date.Format("2006-01-02 15:04:05")

	return orderWithProduct, err
}

// StoreOrder is to store product into database
func (intr Internal) StoreOrder(ctx context.Context, tx *sql.Tx, order Order) (ID int64, err error) {
	var args []interface{}
	query := `INSERT INTO orders  
					(
						order_id_format,
						product_id,
						quantity,
						description,
						date,
						price 
					)
			VALUES (
						?, 
						?, 
						?,
						?, 
						?, 
						?						
					)
			`
	args = append(args,
		order.OrderIDFormat,
		order.ProductID,
		order.Quantity,
		order.Description,
		order.Date,
		order.Price,
	)
	if order.OrderID != 0 {
		query = `UPDATE orders 
				 SET 
						order_id_format = ?,
						product_id = ?,
						quantity = ?,
						description = ?,
						date = ?,
						price = ? 
				WHERE 
						order_id = ?
						 
		`
		args = append(args, order.OrderID)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	if order.OrderID == 0 {
		// no need to check error, since it will be occurred by database incompatibility
		order.OrderID, _ = result.LastInsertId()
	}

	return order.OrderID, err
}
