package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sog01/ijahshop/handler/internal"
	"github.com/sog01/ijahshop/module"
)

// GetProduct is to serve API which get all product
func (h API) GetProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.mod.GetProduct(r.Context())
	if err != nil {
		log.Printf("Error Get Product [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "products", products)
}

// GetDetailProduct is to serve API which get one product
func (h API) GetDetailProduct(w http.ResponseWriter, r *http.Request) {
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

	product, err := h.mod.GetProductByID(r.Context(), ID)
	if err != nil {
		log.Printf("Error Get Product By ID [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "product", product)
}

// StoreProduct is to serve API which store product into database
func (h API) StoreProduct(w http.ResponseWriter, r *http.Request) {
	var reqProduct module.ReqProduct

	// validate request of json
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&reqProduct)
	if err != nil {
		log.Printf("bad request [err = %v, body = %+v]\n", err, r.Body)
		internal.ConstructRespErrorWithDetail(w, http.StatusBadRequest, "Bad Request", map[string]interface{}{
			"description": "Invalid Json Request",
			"info":        "json format : product_id, name, sku, stock",
		})
		return
	}

	reqProduct.ProductID, err = h.mod.StoreProduct(r.Context(), reqProduct)
	if err != nil {
		log.Printf("Error Store product into database [err = %v], [req = %+v]\n", err, reqProduct)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "product", reqProduct)
}

// DeleteProduct is to serve API which delete product
func (h API) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	err = h.mod.DeleteProduct(r.Context(), ID)
	if err != nil {
		log.Printf("Error Get Product By ID [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.ConstructRespSucces(w, "product", map[string]interface{}{"success": "true"})
}

// GetProductCSV is to serve API which get csv file of entity product
func (h API) GetProductCSV(w http.ResponseWriter, r *http.Request) {
	err := h.mod.WriteProductToCSV(r.Context())
	if err != nil {
		log.Printf("Error Get Product CSV [%v]\n", err)
		internal.ConstructRespError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	internal.DownloadFile(w, "Catatan Jumlah Barang.csv")
}
