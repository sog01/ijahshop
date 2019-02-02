package module

import (
	"context"
	"time"

	"github.com/sog01/ijahshop/module/internal"
)

// ReqOrder is entity of inputed order
// use to make request that will be stored into database
type ReqOrder struct {
	internal.Order
	DateRaw string `json:"date_raw"`
}

// ReqFilterOrder is entity to filter query
type ReqFilterOrder struct {
	DateStart time.Time
	DateEnd   time.Time
}

// GetOrderWithProduct is used to get all order with product
func (mod Module) GetOrderWithProduct(ctx context.Context, reqFilter ReqFilterOrder) ([]internal.OrderWithProduct, error) {
	var err error
	ordersWithProduct := []internal.OrderWithProduct{}
	if reqFilter.DateStart != (time.Time{}) && reqFilter.DateEnd != (time.Time{}) {
		ordersWithProduct, err = mod.internal.GetOrderWithProductByDate(ctx, reqFilter.DateStart, reqFilter.DateEnd)
	} else {
		ordersWithProduct, err = mod.internal.GetOrderWithProduct(ctx)
	}

	if err != nil {
		return nil, err
	}

	// calculate total
	for index, orderWithProduct := range ordersWithProduct {
		ordersWithProduct[index].Total = orderWithProduct.Price * int64(orderWithProduct.Quantity)
	}

	return ordersWithProduct, nil
}

// GetOrderWithProductByID is used to get product by ID
func (mod Module) GetOrderWithProductByID(ctx context.Context, ID int64) (internal.OrderWithProduct, error) {
	orderWithProduct, err := mod.internal.GetOrderWithProductByID(ctx, ID)
	if err != nil {
		return internal.OrderWithProduct{}, err
	}

	// calculate total
	orderWithProduct.Total = orderWithProduct.Price * int64(orderWithProduct.Quantity)

	return orderWithProduct, nil
}

// StoreOrder is to store order into database
func (mod Module) StoreOrder(ctx context.Context, reqOrder ReqOrder) (ID int64, err error) {

	// date format: yyyy-MM-dd HH:mm:ss
	date, err := time.Parse("2006-01-02 15:04:05", reqOrder.DateRaw)
	if err != nil {
		return 0, err
	}
	order := internal.Order{
		OrderID:       reqOrder.OrderID,
		OrderIDFormat: reqOrder.OrderIDFormat,
		ProductID:     reqOrder.ProductID,
		Quantity:      reqOrder.Quantity,
		Description:   reqOrder.Description,
		Date:          date,
		Price:         reqOrder.Price,
	}

	db := mod.Storage.DB
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	ID, err = mod.internal.StoreOrder(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return ID, tx.Commit()
}

// WriteOrderToCSV to write order entity to CSV
func (mod Module) WriteOrderToCSV(ctx context.Context) error {
	order, err := mod.GetOrderWithProduct(ctx, ReqFilterOrder{})
	if err != nil {
		return err
	}
	return mod.writeToCSV(ctx, "Catatan Barang Keluar", order)
}
