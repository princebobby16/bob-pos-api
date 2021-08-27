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

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()
	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}
	traceId, productId := getProductIdAndTraceId(headers, r)
	if deleteProductEntryFromDatabase(w, err, productId, transactionId) {
		return
	}
	sendDeleteProductSuccessResponse(w, productId, transactionId, traceId)
}

func getProductIdAndTraceId(headers map[string]string, r *http.Request) (string, string) {
	traceId := headers["trace-id"]
	// Logging the headers
	logger.Logger.Infof("Headers => TraceId: %s", traceId)

	productId := r.URL.Query().Get("product_id")
	logger.Logger.Info(productId)
	return traceId, productId
}

func sendDeleteProductSuccessResponse(w http.ResponseWriter, productId string, transactionId uuid.UUID, traceId string) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        productId,
			UiMessage: "Product Deleted!",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func deleteProductEntryFromDatabase(w http.ResponseWriter, err error, productId string, transactionId uuid.UUID) bool {
	query := `delete from bobpos.products where id=$1`

	_, err = connection.Connection.Exec(query, &productId)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return true
	}
	return false
}
