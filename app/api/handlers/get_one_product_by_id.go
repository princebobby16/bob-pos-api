package handlers

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/db/connection"
	"gitlab.com/pbobby001/bobpos_api/pkg/logger"
	"net/http"
	"time"
)

func GetOneProductById(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}
	//Get the relevant headers
	traceId := headers["trace-id"]
	// Logging the headers
	logger.Logger.Infof("Headers => TraceId: %s", traceId)

	productId := r.URL.Query().Get("product_id")
	logger.Logger.Info(productId)

	product, done := getProductFromDatabase(w, err, productId, transactionId)
	if done {
		return
	}

	sendSuccessResponse(w, product, transactionId, traceId)
}

func getProductFromDatabase(w http.ResponseWriter, err error, productId string, transactionId uuid.UUID) (pkg.Product, bool) {
	query := `select * from bobpos.products where id=$1`
	var product pkg.Product
	err = connection.Connection.QueryRow(query, &productId).Scan(
		&product.Id,
		&product.Name,
		&product.Category,
		&product.Weight,
		&product.CostPrice,
		&product.Tax,
		&product.ProfitMargin,
		&product.Image,
		&product.NumberInStock,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Barcode,
	)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return pkg.Product{}, true
	}
	return product, false
}

func sendSuccessResponse(w http.ResponseWriter, product pkg.Product, transactionId uuid.UUID, traceId string) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pkg.GetOneProductResponse{
		Product: product,
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}
