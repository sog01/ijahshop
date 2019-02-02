package module

import (
	"context"
	"time"

	"github.com/sog01/ijahshop/module/internal"
)

// ReqPurchase is entity of inputed purchase with purchase detail
// use to make request that will be stored into database
type ReqPurchase struct {
	internal.Purchase
	DateRaw     string           `json:"date_raw"`
	PurchaseDtl []ReqPurchaseDtl `json:"purchase_dtl"`
}

// ReqPurchaseDtl is entity of inputed purchase detail
// use to make request that will be stored into database
type ReqPurchaseDtl struct {
	internal.PurchaseDtl
	DateRaw string `json:"date_raw"`
}

// ReqFilterPurchase is entity to filter query
type ReqFilterPurchase struct {
	DateStart time.Time
	DateEnd   time.Time
}

// GetPurchaseWithProduct is used to get all purchase with product
func (mod Module) GetPurchaseWithProduct(ctx context.Context, reqFilter ReqFilterPurchase) ([]internal.PurchaseWithProduct, error) {
	var err error
	purchasesProduct := []internal.PurchaseWithProduct{}

	if reqFilter.DateStart != (time.Time{}) && reqFilter.DateEnd != (time.Time{}) {
		purchasesProduct, err = mod.internal.GetPurchaseWithProductByDate(ctx, reqFilter.DateStart, reqFilter.DateEnd)
	} else {
		purchasesProduct, err = mod.internal.GetPurchaseWithProduct(ctx)
	}

	if err != nil {
		return nil, err
	}

	// calculate total
	for index, purchaseProduct := range purchasesProduct {
		purchasesProduct[index].Total = purchaseProduct.Cost * int64(purchaseProduct.QuantityOrder)
	}

	return purchasesProduct, nil
}

// GetPurchaseWithProductByID is used to get purchase with product by ID
func (mod Module) GetPurchaseWithProductByID(ctx context.Context, ID int64) (internal.PurchaseWithProduct, error) {

	purchaseProduct, err := mod.internal.GetPurchaseWithProductByID(ctx, ID)
	if err != nil {
		return internal.PurchaseWithProduct{}, err
	}

	// calculate total
	purchaseProduct.Total = purchaseProduct.Cost * int64(purchaseProduct.QuantityOrder)

	return purchaseProduct, nil
}

// StorePurchase is to store purchase into database
func (mod Module) StorePurchase(ctx context.Context, reqPurchase ReqPurchase) (ID int64, err error) {

	// date format: yyyy-MM-dd HH:mm:ss
	reqPurchase.Date, err = time.Parse("2006-01-02 15:04:05", reqPurchase.DateRaw)
	if err != nil {
		return 0, err
	}
	purchase := internal.Purchase{
		PurchaseID:       reqPurchase.PurchaseID,
		ProductID:        reqPurchase.ProductID,
		QuantityOrder:    reqPurchase.QuantityOrder,
		QuantityAccepted: reqPurchase.QuantityAccepted,
		Description:      reqPurchase.Description,
		InvoiceNumber:    reqPurchase.InvoiceNumber,
		Cost:             reqPurchase.Cost,
		Date:             reqPurchase.Date,
		IsFinish:         reqPurchase.IsFinish,
	}

	db := mod.Storage.DB
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	purchaseID, err := mod.internal.StorePurchase(ctx, tx, purchase)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, reqPurchaseDtl := range reqPurchase.PurchaseDtl {
		// date format: yyyy-MM-dd HH:mm:ss
		reqPurchaseDtl.Date, err = time.Parse("2006-01-02 15:04:05", reqPurchaseDtl.DateRaw)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		purchaseDtl := internal.PurchaseDtl{
			PurchaseDtlID: reqPurchaseDtl.PurchaseDtlID,
			PurchaseID:    purchaseID,
			Quantity:      reqPurchaseDtl.Quantity,
			Description:   reqPurchaseDtl.Description,
			Date:          reqPurchaseDtl.Date,
		}
		_, err = mod.internal.StorePurchaseDtl(ctx, tx, purchaseDtl)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return purchaseID, tx.Commit()
}

// WritePurchaseToCSV to write purchase entity to CSV
func (mod Module) WritePurchaseToCSV(ctx context.Context) error {
	purchase, err := mod.GetPurchaseWithProduct(ctx, ReqFilterPurchase{})
	if err != nil {
		return err
	}
	return mod.writeToCSV(ctx, "Catatan Barang Masuk", purchase)
}
