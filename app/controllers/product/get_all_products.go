package product

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/db"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/logs"
	"net/http"
	"time"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	logs.Logger.Info("in get")
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	logs.Logger.Infof("Headers => TraceId: %s", traceId)

	query := `select * from bobpos.products limit 2000`

	rows, err := db.Connection.Query(query)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	var products []pkg.Product
	for rows.Next() {
		var product pkg.Product
		err = rows.Scan(
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
		)
		if err != nil {
			pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
			return
		}

		products = append(products, product)
		product.Image = []byte{}
		logs.Logger.Info(product)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pkg.StandardGetAllProductsResponse{
		Data: products,
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}
