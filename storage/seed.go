package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

// DataSeed is entity of data that want to seed into database
type DataSeed struct {
	Table      string
	DataColumn []dataColumn
}

type dataColumn struct {
	Data   interface{}
	Column string
}

// Seed to seed data into database
func (s Storage) Seed(data []DataSeed) error {
	db := s.DB
	for _, row := range data {
		var columnStr, args []string
		var data []interface{}
		query := "INSERT INTO " + row.Table + "(%s) VALUES (%s)"
		for _, val := range row.DataColumn {
			// construct column dan data string
			columnStr = append(columnStr, val.Column)
			data = append(data, val.Data)
			args = append(args, "?")
		}
		query = fmt.Sprintf(query, strings.Join(columnStr, ","), strings.Join(args, ","))
		_, err := db.Exec(db.Rebind(query), data...)
		if err != nil {
			log.Printf("Failed insert into DB [err = %v] [query = %s] ", err, query)
			return err
		}

	}
	return nil
}

// SeedDummyData to seed dummy data into database
func (s Storage) SeedDummyData() error {
	data := []DataSeed{
		DataSeed{
			Table: "product",
			DataColumn: []dataColumn{
				dataColumn{
					Data:   "BULUGUL MARAM",
					Column: "name",
				},
				dataColumn{
					Data:   "SSI-D00791077-MM-BM",
					Column: "sku",
				},
				dataColumn{
					Data:   10,
					Column: "stock",
				},
			},
		},
		DataSeed{
			Table: "product",
			DataColumn: []dataColumn{
				dataColumn{
					Data:   "RIYADUS SHALIHIN",
					Column: "name",
				},
				dataColumn{
					Data:   "SSI-D00791077-MM-RS",
					Column: "sku",
				},
				dataColumn{
					Data:   20,
					Column: "stock",
				},
			},
		},
	}

	return s.Seed(data)
}

// SeedProductFromEXCEL to seed data from excel file
func (s Storage) SeedProductFromEXCEL(filePath string) error {
	// mapping sheet from excel into table in db
	mapTable := make(map[string]string)
	mapTable["Catatan Jumlah Barang"] = "product"
	mapTable["Catatan Barang Masuk"] = "purchase"
	mapTable["Catatan Barang Keluar"] = "orders"

	// mapping column from excel into field in db
	mapColumn := make(map[string]string)
	mapColumn["SKU"] = "sku"
	mapColumn["Nama Item"] = "name"
	mapColumn["Jumlah Sekarang"] = "stock"
	mapColumn["Jumlah Pemesanan"] = "quantity_order"
	mapColumn["Jumlah Diterima"] = "quantity_accepted"
	mapColumn["Harga Beli"] = "cost"
	mapColumn["Nomer Kwitansi"] = "invoice_number"
	mapColumn["Catatan"] = "description"
	mapColumn["Waktu"] = "date"
	mapColumn["Jumlah Keluar"] = "quantity"
	mapColumn["Harga Jual"] = "price"

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return err
	}
	for _, sheet := range xlFile.Sheets {
		var (
			datas  []DataSeed
			column []string
		)
		for index, row := range sheet.Rows {
			var data DataSeed
			data.Table = mapTable[sheet.Name]
			if data.Table == "" {
				continue
			}
			for key, cell := range row.Cells {
				var (
					columnData dataColumn
					columnName string
				)
				text := cell.String()
				if index == 0 {
					text = strings.Trim(text, " ")
					column = append(column, mapColumn[text])
				} else {
					if text == "" || column[key] == "" {
						continue
					}
					columnName = column[key]
					if column[key] == "date" {
						// date format: yyyy-MM-dd HH:mm:ss
						date, err := time.Parse("2006-01-02 15:04:05", text)
						text = date.String()
						if err != nil {
							// set date for today
							text = time.Now().Format("2006-01-02 15:04:05")
						}
					} else if column[key] == "price" || column[key] == "cost" {
						// normalize data which has value format : "Rp74000
						text = strings.Replace(text, "Rp", "", -1)
						text = strings.Replace(text, ",", "", -1)
					} else if data.Table != "product" {
						if columnName == "sku" {
							productID, err := s.skuToProductID(text)
							if err != nil {
								log.Println(data.Table, text, err)
							}
							columnName = "product_id"
							text = fmt.Sprintf("%d", productID)
						}
					}

					if data.Table == "orders" {
						if columnName == "description" {
							var orderIDFormat string
							bufferSplit := strings.Split(text, " ")
							if len(bufferSplit) == 2 {
								orderIDFormat = bufferSplit[1]
							}
							columnData.Column = "order_id_format"
							columnData.Data = orderIDFormat
							data.DataColumn = append(data.DataColumn, columnData)
						}
					}
					columnData.Column = columnName
					columnData.Data = text
				}
				data.DataColumn = append(data.DataColumn, columnData)
			}
			// since index == 0 is a column, so no need to collect data
			if index > 0 {
				if len(data.DataColumn) > 0 {
					datas = append(datas, data)
				}
			}
		}
		err := s.Seed(datas)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) skuToProductID(sku string) (int64, error) {
	var productID int64
	db := s.DB
	query := "SELECT product_id FROM product WHERE sku = ?"
	row := db.QueryRow(db.Rebind(query), sku)
	err := row.Scan(&productID)

	return productID, err
}
