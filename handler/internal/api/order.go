package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/sog01/ijahshop/handler/internal"
	"github.com/sog01/ijahshop/module"
)

// GetOrder is to serve API which get all Order
func (h API) GetOrder(w http.ResponseWriter, r *http.Request) {
	orders, err := h.mod.GetOrderWithProduct(r.Context(), module.ReqFilterOrder{})
	if err != nil {
		log.Printf("Error Get Order [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "orders", orders)
}

// GetOrderByDate is to serve API which get all Order by date
func (h API) GetOrderByDate(w http.ResponseWriter, r *http.Request) {
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

	ordersWithProduct, err := h.mod.GetOrderWithProduct(r.Context(), module.ReqFilterOrder{
		DateStart: dateStart,
		DateEnd:   dateEnd,
	})
	if err != nil {
		log.Printf("Error Get Order [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "orders", ordersWithProduct)
}

// GetDetailOrder is to serve API which get one Order
func (h API) GetDetailOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	IDstr := vars["id"]

	// validate request
	ID, err := strconv.ParseInt(IDstr, 10, 64)
	if err != nil {
		log.Printf("Bad Request Form ID [%v]\n", err)
		internal.ConstructRespErrorWithDetail(w, http.StatusBadRequest, "Bad Request", map[string]interface{}{
			"description": "Invalid id",
		})
		return
	}

	order, err := h.mod.GetOrderWithProductByID(r.Context(), ID)
	if err != nil {
		log.Printf("Error Get Order By ID [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "order", order)
}

// StoreOrder is to serve API which store order into database
func (h API) StoreOrder(w http.ResponseWriter, r *http.Request) {
	var reqOrder module.ReqOrder

	// validate request of json
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&reqOrder)
	if err != nil {
		log.Printf("bad request [err = %v, body = %+v]\n", err, r.Body)
		internal.ConstructRespErrorWithDetail(w, http.StatusBadRequest, "Bad Request", map[string]interface{}{
			"description": "Invalid Json Request",
		})
		return
	}

	reqOrder.OrderID, err = h.mod.StoreOrder(r.Context(), reqOrder)
	if err != nil {
		log.Printf("Error Store order into database [err = %v], [req = %+v]\n", err, reqOrder)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "order", reqOrder)
}

// GetOrderCSV is to serve API which get csv file of entity order
func (h API) GetOrderCSV(w http.ResponseWriter, r *http.Request) {
	err := h.mod.WriteOrderToCSV(r.Context())
	if err != nil {
		log.Printf("Error Get Order CSV [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.DownloadFile(w, "Catatan Barang Keluar.csv")
}
