package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sog01/ijahshop/config"
	"github.com/sog01/ijahshop/handler"
	"github.com/sog01/ijahshop/module"
	"github.com/sog01/ijahshop/storage"
	"github.com/sog01/ijahshop/template"
)

func main() {
	configPath := "files/config.ini"
	conf, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Failed create config instance [%v]\n", err)
	}

	dataSource := conf.Storage["sqlite"].Host
	storageDB, err := storage.New(dataSource)
	if err != nil {
		log.Fatalf("Failed create storage instance [%v]\n", err)
	}

	// migration script
	// this script is optional, you can remove it
	// if all table already set up
	err = storageDB.Migrate()
	if err != nil {
		log.Printf("Failed to migrate table [%v]\n", err)
	}

	module := module.New(storageDB)
	if err != nil {
		log.Fatalf("Failed create module instance [%v]\n", err)
	}

	tmpl, err := template.New()
	if err != nil {
		log.Fatalf("Failed create template instance [%v]\n", err)
	}

	handlr := handler.New(module, tmpl)

	r := mux.NewRouter()

	// handle static file
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", handlr.StaticJS()))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", handlr.StaticCSS()))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", handlr.StaticImages()))
	r.PathPrefix("/scripts/").Handler(http.StripPrefix("/scripts/", handlr.StaticScript()))
	r.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", handlr.StaticStyles()))
	r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", handlr.StaticData()))

	// handle API
	{
		// serve product request
		r.HandleFunc("/inventory/product", handlr.API.GetProduct).Methods("GET")
		r.HandleFunc("/inventory/product", handlr.API.StoreProduct).Methods("POST")
		r.HandleFunc("/inventory/product/{id:[0-9]+}", handlr.API.GetDetailProduct).Methods("GET")
		r.HandleFunc("/inventory/product/{id:[0-9]+}", handlr.API.DeleteProduct).Methods("DELETE")
	}

	{
		// serve purchase request
		r.HandleFunc("/inventory/purchase", handlr.API.GetPurchase).Methods("GET")
		r.HandleFunc("/inventory/purchase/{date_start}/{date_end}", handlr.API.GetPurchaseByDate).Methods("GET")
		r.HandleFunc("/inventory/purchase/{id}", handlr.API.GetDetailPurchase).Methods("GET")
		r.HandleFunc("/inventory/purchase", handlr.API.StorePurchase).Methods("POST")
	}

	{
		// serve order request
		r.HandleFunc("/inventory/order", handlr.API.GetOrder).Methods("GET")
		r.HandleFunc("/inventory/order/{date_start}/{date_end}", handlr.API.GetOrderByDate).Methods("GET")
		r.HandleFunc("/inventory/order/{id}", handlr.API.GetDetailOrder).Methods("GET")
		r.HandleFunc("/inventory/order", handlr.API.StoreOrder).Methods("POST")
	}

	{
		// serve report request
		r.HandleFunc("/inventory/report/product", handlr.API.GetProductReport).Methods("GET")
		r.HandleFunc("/inventory/report/order/{date_start}/{date_end}", handlr.API.GetOrderReport).Methods("GET")
	}

	{
		// serve export request
		r.HandleFunc("/inventory/export/product", handlr.API.GetProductCSV).Methods("GET")
		r.HandleFunc("/inventory/export/purchase", handlr.API.GetPurchaseCSV).Methods("GET")
		r.HandleFunc("/inventory/export/order", handlr.API.GetOrderCSV).Methods("GET")
		r.HandleFunc("/inventory/export/report_product", handlr.API.GetProductReportCSV).Methods("GET")
		r.HandleFunc("/inventory/export/report_order/{date_start}/{date_end}", handlr.API.GetOrderReportCSV).Methods("GET")
	}

	// optional task
	{
		// import excel
		r.HandleFunc("/inventory/import", handlr.API.ImportExcelFile).Methods("POST")

		// frontend web
		r.HandleFunc("/", handlr.Product).Methods("GET")
		r.HandleFunc("/purchase", handlr.Purchase).Methods("GET")
		r.HandleFunc("/orders", handlr.Orders).Methods("GET")
		r.HandleFunc("/product/report", handlr.ProductReport).Methods("GET")
		r.HandleFunc("/orders/report", handlr.OrderReport).Methods("GET")
	}

	http.ListenAndServe(":8080", r)
}
