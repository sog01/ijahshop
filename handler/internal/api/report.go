package internal

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sog01/ijahshop/handler/internal"
	"github.com/sog01/ijahshop/module"
)

// GetProductReport is to serve API which get all productAvgValue
func (h API) GetProductReport(w http.ResponseWriter, r *http.Request) {
	productAvgValuesWithProduct, err := h.mod.GetProductAvgValue(r.Context())
	if err != nil {
		log.Printf("Error Get productAvgValue [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// mapping summary which will be shown as meta data
	summary := map[string]interface{}{
		"summary": productAvgValuesWithProduct.Summary,
	}

	internal.ConstructRespSuccesWithMeta(w,
		"product_report",
		productAvgValuesWithProduct.ProductAvgValue,
		summary,
	)
}

// GetOrderReport is to serve API which get all order with average value by date
func (h API) GetOrderReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dateStartStr := vars["date_start"]
	dateStartEnd := vars["date_end"]

	// sanitize request
	// date format: yyyy-MM-dd
	dateStart, err := time.Parse("2006-01-02", dateStartStr)
	if err != nil {
		log.Printf("Bad Request date start [%v]\n", err)
		internal.ConstructRespErrorWithDetail(w, http.StatusInternalServerError, "Bad Request", map[string]interface{}{
			"description": "Invalid date format",
		})
		return
	}

	dateEnd, err := time.Parse("2006-01-02", dateStartEnd)
	if err != nil {
		log.Printf("Bad Request date end [%v]\n", err)
		internal.ConstructRespErrorWithDetail(w, http.StatusInternalServerError, "Bad Request", map[string]interface{}{
			"description": "Invalid date format",
		})
		return
	}

	orderWithProductValueWithSummary, err := h.mod.GetOrderWithProductAvgValueByDate(r.Context(), module.ReqFilterOrder{
		DateStart: dateStart,
		DateEnd:   dateEnd,
	})

	if err != nil {
		log.Printf("Error Get Order [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// mapping summary which will be shown as meta data
	summary := map[string]interface{}{
		"summary": orderWithProductValueWithSummary.Summary,
	}

	internal.ConstructRespSuccesWithMeta(w,
		"order_report",
		orderWithProductValueWithSummary.OrderWithProductValue,
		summary,
	)
}

// GetProductReportCSV is to serve API which get csv file of entity product report
func (h API) GetProductReportCSV(w http.ResponseWriter, r *http.Request) {
	err := h.mod.WriteProductReportToCSV(r.Context())
	if err != nil {
		log.Printf("Error Get Product Report CSV [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.DownloadFile(w, "Laporan Nilai Barang.csv")
}

// GetOrderReportCSV is to serve API which get csv file of entity order report
func (h API) GetOrderReportCSV(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dateStartStr := vars["date_start"]
	dateStartEnd := vars["date_end"]

	// sanitize request
	// date format: yyyy-MM-dd
	dateStart, err := time.Parse("2006-01-02", dateStartStr)
	if err != nil {
		log.Printf("Bad Request date start [%v]\n", err)
		internal.ConstructRespErrorWithDetail(w, http.StatusInternalServerError, "Bad Request", map[string]interface{}{
			"description": "Invalid date format",
		})
		return
	}

	dateEnd, err := time.Parse("2006-01-02", dateStartEnd)
	if err != nil {
		log.Printf("Bad Request date end [%v]\n", err)
		internal.ConstructRespErrorWithDetail(w, http.StatusInternalServerError, "Bad Request", map[string]interface{}{
			"description": "Invalid date format",
		})
		return
	}

	err = h.mod.WriteOrderReportToCSV(r.Context(), module.ReqFilterOrder{
		DateStart: dateStart,
		DateEnd:   dateEnd,
	})
	if err != nil {
		log.Printf("Error Get Order Report CSV [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.DownloadFile(w, "Laporan Penjualan.csv")
}

// ImportExcelFile is to serve API which import csv into database
func (h API) ImportExcelFile(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error Bad Request [%v]\n", err)
		internal.ConstructRespError(w, http.StatusBadRequest, "Bad request")
		return
	}

	defer file.Close()

	f, err := os.OpenFile("files/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error Failed to open files [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer f.Close()
	io.Copy(f, file)

	err = h.mod.ImportExcelToDB(header.Filename)
	if err != nil {
		log.Printf("Error Failed to import files [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "success", true)

}
