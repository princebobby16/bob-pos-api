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

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	transactionId, err, traceId, done := generateTransactionIdAndExtractTraceId(w, r)
	if done {
		return
	}

	categories, done2 := getAllCategoriesFromTheDatabase(w, err, transactionId, traceId)
	if done2 {
		return
	}
	logger.Logger.Info(categories)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pkg.GetAllCategoryResponse{
		Data: categories,
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func getAllCategoriesFromTheDatabase(w http.ResponseWriter, err error, transactionId uuid.UUID, traceId string) ([]pkg.ProductCategory, bool) {
	query := `select * from bobpos.product_category`

	rows, err := connection.Connection.Query(query)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
		return nil, true
	}

	var categories []pkg.ProductCategory
	for rows.Next() {
		var category pkg.ProductCategory
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
			return nil, true
		}
		categories = append(categories, category)
	}
	return categories, false
}

func generateTransactionIdAndExtractTraceId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error, string, bool) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", true
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	logger.Logger.Infof("Headers => TraceId: %s", traceId)
	return transactionId, err, traceId, false
}
