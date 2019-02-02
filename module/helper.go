package module

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/sog01/ijahshop/module/internal"
)

func (mod Module) writeToCSV(ctx context.Context, filename string, data interface{}) error {
	file, err := os.Create("files/data/" + filename + ".csv")
	defer file.Close()
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	dataRow, ok := data.([][]string)
	if !ok {
		dataRow = mod.entityIntoArrayString(data)
	}

	for _, value := range dataRow {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}

	return nil

}

func (mod Module) entityIntoArrayString(entity interface{}) [][]string {
	var arrString [][]string

	mapType := map[string]string{
		"internal.Product":                    "product",
		"internal.PurchaseWithProduct":        "purchase",
		"internal.OrderWithProduct":           "order",
		"internal.ProductAvgValue":            "report_product",
		"module.OrderWithProductValue":        "report_order",
		"module.SummaryAvgValue":              "report_product_summary",
		"module.SummaryOrderWithProductValue": "report_order_summary",
	}

	// internal function
	extractToRow := func(e reflect.Value, mapper map[string]string, indexRow int) []string {
		var row []string
		for j := 0; j < e.NumField(); j++ {
			varName := e.Type().Field(j).Name
			varValue := e.Field(j).Interface()

			// retrive column in row index - 0
			if indexRow == 0 {
				if mapper[varName] != "" {
					row = append(row, mapper[varName])
				}
			} else {
				if mapper[varName] != "" {
					row = append(row, fmt.Sprintf("%v", varValue))
				}
			}
		}
		return row
	}

	object := reflect.ValueOf(entity)

	if object.Kind() == reflect.Slice {
		for i := 0; i < object.Len(); i++ {
			var (
				column []string
				row    []string
			)
			obj := object.Index(i)
			typeObj := obj.Type().String()
			switch mapType[typeObj] {
			case "product":
				product, ok := obj.Interface().(internal.Product)
				if !ok {
					return nil
				}

				e := reflect.ValueOf(&product).Elem()
				mapper := map[string]string{
					"Name":  "Nama Item",
					"Sku":   "SKU",
					"Stock": "Jumlah Sekarang",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

			case "purchase":
				purchase, ok := obj.Interface().(internal.PurchaseWithProduct)
				if !ok {
					return nil
				}
				e := reflect.ValueOf(&purchase.Purchase).Elem()
				mapper := map[string]string{
					"DateStr":       "Waktu",
					"OrderQuantity": "Jumlah Pemesanan",
					"OrderAccepted": "Jumlah Diterima",
					"Cost":          "Harga Beli",
					"Total":         "Total",
					"InvoiceNumber": "Nomer Kuitansi",
					"Description":   "Catatan",
				}
				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

				e = reflect.ValueOf(&purchase.Product).Elem()
				mapper = map[string]string{
					"Name": "Nama Barang",
					"Sku":  "SKU",
				}
				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

			case "order":
				order, ok := obj.Interface().(internal.OrderWithProduct)
				if !ok {
					return nil
				}
				e := reflect.ValueOf(&order.Order).Elem()
				mapper := map[string]string{
					"Date":        "Waktu",
					"Quantity":    "Jumlah Keluar",
					"Price":       "Harga Jual",
					"Total":       "Total",
					"Description": "Catatan",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

				e = reflect.ValueOf(&order.Product).Elem()
				mapper = map[string]string{
					"Name": "Nama Barang",
					"Sku":  "SKU",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

			case "report_product":
				reportProduct, ok := obj.Interface().(internal.ProductAvgValue)
				if !ok {
					return nil
				}
				e := reflect.ValueOf(&reportProduct.Product).Elem()
				mapper := map[string]string{
					"Name":  "Nama Item",
					"Sku":   "SKU",
					"Stock": "Jumlah",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

				e = reflect.ValueOf(&reportProduct).Elem()
				mapper = map[string]string{
					"AverageCost": "Rata-Rata Harga Beli",
					"Total":       "Total",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

			case "report_order":
				reportOrder, ok := obj.Interface().(OrderWithProductValue)
				if !ok {
					return nil
				}
				e := reflect.ValueOf(&reportOrder.Order).Elem()
				mapper := map[string]string{
					"OrderIDFormat": "ID Pesanan",
					"Date":          "Waktu",
					"Quantity":      "Jumlah",
					"Price":         "Harga Jual",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

				e = reflect.ValueOf(&reportOrder).Elem()
				mapper = map[string]string{
					"AverageCost": "Harga Beli",
					"Total":       "Total",
					"Profit":      "Laba",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

				e = reflect.ValueOf(&reportOrder.Product).Elem()
				mapper = map[string]string{
					"Name": "Nama Barang",
					"Sku":  "SKU",
				}

				if i == 0 {
					column = append(column, extractToRow(e, mapper, 0)...)
				}
				row = append(row, extractToRow(e, mapper, 1)...)

			default:
				return nil
			}
			if i == 0 {
				arrString = append(arrString, column)
			}
			arrString = append(arrString, row)
		}

	} else {
		typeObj := object.Type().String()
		switch mapType[typeObj] {
		case "report_product_summary":
			summary, ok := object.Interface().(SummaryAvgValue)
			if !ok {
				return nil
			}

			var rows [][]string
			rows = append(rows, []string{"Tanggal Cetak : " + summary.DatePrint})
			rows = append(rows, []string{fmt.Sprintf("Jumlah SKU : %d", summary.TotalSku)})
			rows = append(rows, []string{fmt.Sprintf("Jumlah Total Barang : %d", summary.TotalProduct)})
			rows = append(rows, []string{fmt.Sprintf("Total Nilai : %d", summary.TotalValue)})

			arrString = append(arrString, rows...)
		case "report_order_summary":
			summary, ok := object.Interface().(SummaryOrderWithProductValue)
			if !ok {
				return nil
			}

			var rows [][]string
			rows = append(rows, []string{"Tanggal Cetak : " + summary.DatePrint})
			rows = append(rows, []string{"Tanggal : " + summary.Date})
			rows = append(rows, []string{fmt.Sprintf("Total Omzet : %d", summary.TotalPrice)})
			rows = append(rows, []string{fmt.Sprintf("Laba Kotor : %d", summary.TotalProfit)})
			rows = append(rows, []string{fmt.Sprintf("Total Penjualan : %d", summary.TotalSold)})
			rows = append(rows, []string{fmt.Sprintf("Total Barang : %d", summary.TotalItem)})

			arrString = append(arrString, rows...)
		}

	}

	return arrString
}
