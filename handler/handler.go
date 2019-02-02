package handler

import (
	"net/http"
	"time"

	"github.com/sog01/ijahshop/handler/internal/api"
	"github.com/sog01/ijahshop/module"
	tmplInt "github.com/sog01/ijahshop/template"
)

// Handler is main entity of package handler
type Handler struct {
	tmpl tmplInt.Template
	mod  module.Module
	API  internal.API
}

// New to create new instance of handler main entity
func New(mod module.Module, tmpl tmplInt.Template) Handler {
	handler := Handler{
		tmpl: tmpl,
		mod:  mod,
		API:  internal.New(mod),
	}
	return handler
}

// Product is http func that handle index page
func (h Handler) Product(w http.ResponseWriter, r *http.Request) {
	finalTemplate := h.tmpl["product"]

	product, _ := h.mod.GetProduct(r.Context())

	finalTemplate.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Data": product,
	})
}

// Purchase is http func that handle Purchase page
func (h Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	finalTemplate := h.tmpl["purchase"]

	purchase, _ := h.mod.GetPurchaseWithProduct(r.Context(), module.ReqFilterPurchase{})

	finalTemplate.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Data": purchase,
	})
}

// Orders is http func that handle Orders page
func (h Handler) Orders(w http.ResponseWriter, r *http.Request) {
	finalTemplate := h.tmpl["orders"]

	orders, _ := h.mod.GetOrderWithProduct(r.Context(), module.ReqFilterOrder{})

	finalTemplate.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Data": orders,
	})
}

// ProductReport is http func that handle ProductReport page
func (h Handler) ProductReport(w http.ResponseWriter, r *http.Request) {
	finalTemplate := h.tmpl["product_report"]

	productReport, _ := h.mod.GetProductAvgValue(r.Context())

	finalTemplate.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Data":    productReport.ProductAvgValue,
		"Summary": productReport.Summary,
	})
}

// OrderReport is http func that handle OrderReport page
func (h Handler) OrderReport(w http.ResponseWriter, r *http.Request) {
	finalTemplate := h.tmpl["orders_report"]

	dateStart, _ := time.Parse("2006-01-02", "2018-01-01")
	dateEnd, _ := time.Parse("2006-01-02", "2018-01-10")

	orderReport, _ := h.mod.GetOrderWithProductAvgValueByDate(r.Context(), module.ReqFilterOrder{
		DateStart: dateStart,
		DateEnd:   dateEnd,
	})

	finalTemplate.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Data":    orderReport.OrderWithProductValue,
		"Summary": orderReport.Summary,
	})
}
