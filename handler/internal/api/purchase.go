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

// GetPurchase is to serve API which get all Purchase
func (h API) GetPurchase(w http.ResponseWriter, r *http.Request) {
	purchasesWithProduct, err := h.mod.GetPurchaseWithProduct(r.Context(), module.ReqFilterPurchase{})
	if err != nil {
		log.Printf("Error Get Purchase [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "purchases", purchasesWithProduct)
}

// GetPurchaseByDate is to serve API which get all Purchase by date
func (h API) GetPurchaseByDate(w http.ResponseWriter, r *http.Request) {
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

	purchasesWithProduct, err := h.mod.GetPurchaseWithProduct(r.Context(), module.ReqFilterPurchase{
		DateStart: dateStart,
		DateEnd:   dateEnd,
	})
	if err != nil {
		log.Printf("Error Get Purchase [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "purchases", purchasesWithProduct)
}

// GetDetailPurchase is to serve API which get one Purchase
func (h API) GetDetailPurchase(w http.ResponseWriter, r *http.Request) {
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

	purchaseWithProduct, err := h.mod.GetPurchaseWithProductByID(r.Context(), ID)
	if err != nil {
		log.Printf("Error Get Purchase By ID [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "purchase", purchaseWithProduct)
}

// StorePurchase is to serve API which store purchase into database
func (h API) StorePurchase(w http.ResponseWriter, r *http.Request) {
	var reqPurchase module.ReqPurchase

	// validate request of json
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&reqPurchase)
	if err != nil {
		log.Printf("bad request [err = %v, body = %+v]\n", err, r.Body)
		internal.ConstructRespErrorWithDetail(w, http.StatusBadRequest, "Bad Request", map[string]interface{}{
			"description": "Invalid Json Request",
		})
		return
	}

	reqPurchase.PurchaseID, err = h.mod.StorePurchase(r.Context(), reqPurchase)
	if err != nil {
		log.Printf("Error Store purchase into database [err = %v], [req = %+v]\n", err, reqPurchase)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "purchase", reqPurchase)
}

// GetPurchaseCSV is to serve API which get csv file of purchase entity
func (h API) GetPurchaseCSV(w http.ResponseWriter, r *http.Request) {
	err := h.mod.WritePurchaseToCSV(r.Context())
	if err != nil {
		log.Printf("Error Get Purchase CSV [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.DownloadFile(w, "Catatan Barang Masuk.csv")
}
