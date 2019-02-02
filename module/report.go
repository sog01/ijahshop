package module

import (
	"context"
	"time"

	"github.com/sog01/ijahshop/module/internal"
)

// ProductAvgValueWithSummary is entity of product value with summary
type ProductAvgValueWithSummary struct {
	ProductAvgValue []internal.ProductAvgValue
	Summary         SummaryAvgValue `json:"summary"`
}

// OrderWithProductValue is entity of order with product average value
type OrderWithProductValue struct {
	internal.OrderWithProduct
	AverageCost int   `db:"average_cost" json:"average_cost"`
	Total       int64 `json:"total"`
	Profit      int64 `json:"profit"`
}

// OrderWithProductValueWithSummary is entity of order with product average value with summary
type OrderWithProductValueWithSummary struct {
	OrderWithProductValue []OrderWithProductValue
	Summary               SummaryOrderWithProductValue `json:"summary"`
}

// SummaryAvgValue is summary of average value product
// which consist of few elements
type SummaryAvgValue struct {
	DatePrint    string `json:"date_print"`
	TotalSku     int    `json:"total_sku"`
	TotalProduct int    `json:"total_product"`
	TotalValue   int    `json:"total_value"`
}

// SummaryOrderWithProductValue is summary of order with product value
// which consist of few elements
type SummaryOrderWithProductValue struct {
	DatePrint   string `json:"date_print"`
	Date        string `json:"date"`
	TotalPrice  int64  `json:"total_price"`
	TotalProfit int64  `json:"total_profit"`
	TotalSold   int    `json:"total_sold"`
	TotalItem   int    `json:"total_item"`
}

// GetProductAvgValue is used to get all product with average value
func (mod Module) GetProductAvgValue(ctx context.Context) (ProductAvgValueWithSummary, error) {
	var productAvgValueWithSummary ProductAvgValueWithSummary

	productsAvgValue, err := mod.internal.GetProductAvgValue(ctx)
	if err != nil {
		return ProductAvgValueWithSummary{}, err
	}

	productAvgValueWithSummary.ProductAvgValue = productsAvgValue

	// calculate total and summary
	// date format: yyyy-MM-dd HH:mm:ss
	productAvgValueWithSummary.Summary.DatePrint = time.Now().Format("2006-01-02 15:04:05")
	for index, productAvgValue := range productAvgValueWithSummary.ProductAvgValue {
		productAvgValueWithSummary.ProductAvgValue[index].Total = productAvgValue.AverageCost * productAvgValue.Stock
		productAvgValueWithSummary.Summary.TotalSku++
		productAvgValueWithSummary.Summary.TotalProduct += productAvgValue.Stock
		productAvgValueWithSummary.Summary.TotalValue += productAvgValueWithSummary.ProductAvgValue[index].Total
	}

	return productAvgValueWithSummary, nil
}

// GetOrderWithProductAvgValueByDate is used to get order with product average value
func (mod Module) GetOrderWithProductAvgValueByDate(ctx context.Context, reqFilter ReqFilterOrder) (OrderWithProductValueWithSummary, error) {
	var (
		orderWithProductValueWithSummary OrderWithProductValueWithSummary
		ordersWithProductValue           []OrderWithProductValue
		summary                          SummaryOrderWithProductValue
	)

	ordersWithProduct, err := mod.internal.GetOrderWithProductByDate(ctx, reqFilter.DateStart, reqFilter.DateEnd)
	if err != nil {
		return OrderWithProductValueWithSummary{}, err
	}

	// date format: yyyy-MM-dd HH:mm:ss
	summary.DatePrint = time.Now().Format("2006-01-02 15:04:05")

	// date has format: date start - date end
	summary.Date = reqFilter.DateStart.Format("2006-01-02 15:04:05") + "-" + reqFilter.DateEnd.Format("2006-01-02 15:04:05")

	for _, orderWithProduct := range ordersWithProduct {
		var orderWithProductValue OrderWithProductValue
		productAvgValue, err := mod.internal.GetProductAvgValueByProductID(ctx, orderWithProduct.ProductID)
		if err != nil {
			return OrderWithProductValueWithSummary{}, err
		}

		orderWithProductValue.OrderID = orderWithProduct.OrderID
		orderWithProductValue.OrderIDFormat = orderWithProduct.OrderIDFormat
		orderWithProductValue.ProductID = orderWithProduct.ProductID
		orderWithProductValue.Quantity = orderWithProduct.Quantity
		orderWithProductValue.Description = orderWithProduct.Description
		orderWithProductValue.Date = orderWithProduct.Date
		orderWithProductValue.DateStr = orderWithProduct.DateStr
		orderWithProductValue.Price = orderWithProduct.Price
		orderWithProductValue.Product = orderWithProduct.Product
		orderWithProductValue.Total = orderWithProductValue.Price * int64(orderWithProductValue.Quantity)
		orderWithProductValue.AverageCost = productAvgValue.AverageCost
		orderWithProductValue.Profit = orderWithProductValue.Total - int64(orderWithProductValue.AverageCost)

		summary.TotalPrice += orderWithProductValue.Total
		summary.TotalProfit += orderWithProductValue.Profit
		summary.TotalItem += orderWithProductValue.Quantity
		summary.TotalSold++

		ordersWithProductValue = append(ordersWithProductValue, orderWithProductValue)
	}

	orderWithProductValueWithSummary.OrderWithProductValue = ordersWithProductValue
	orderWithProductValueWithSummary.Summary = summary

	return orderWithProductValueWithSummary, nil

}

// WriteProductReportToCSV to write product report entity to CSV
func (mod Module) WriteProductReportToCSV(ctx context.Context) error {
	productReport, err := mod.GetProductAvgValue(ctx)
	if err != nil {
		return err
	}

	productAvgValue := mod.entityIntoArrayString(productReport.ProductAvgValue)
	summary := mod.entityIntoArrayString(productReport.Summary)

	row := summary

	// add enter
	row = append(row, []string{})

	row = append(row, productAvgValue...)

	return mod.writeToCSV(ctx, "Laporan Nilai Barang", row)
}

// WriteOrderReportToCSV to write order report entity to CSV
func (mod Module) WriteOrderReportToCSV(ctx context.Context, reqFilter ReqFilterOrder) error {
	orderReport, err := mod.GetOrderWithProductAvgValueByDate(ctx, reqFilter)
	if err != nil {
		return err
	}

	orderWithProductValue := mod.entityIntoArrayString(orderReport.OrderWithProductValue)
	summary := mod.entityIntoArrayString(orderReport.Summary)

	row := summary

	// add enter
	row = append(row, []string{})

	row = append(row, orderWithProductValue...)

	return mod.writeToCSV(ctx, "Laporan Penjualan", row)
}
