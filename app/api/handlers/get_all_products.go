package handlers

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/db/connection"
	"log"
	"net/http"
	"time"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}
	//Get the relevant headers

	traceId := headers["trace-id"]

	products, done := getAllProductsFromDatabase(w, err, transactionId, traceId)
	if done {
		return
	}

	sendGetAllProductsResponse(w, products, transactionId, traceId)
}

func sendGetAllProductsResponse(w http.ResponseWriter, products []pkg.Product, transactionId uuid.UUID, traceId string) {
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

func getAllProductsFromDatabase(w http.ResponseWriter, err error, transactionId uuid.UUID, traceId string) ([]pkg.Product, bool) {
	log.Println("TraceId: ", traceId)

	query := `select * from bobpos.products limit 2000`

	rows, err := connection.Connection.Query(query)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return nil, true
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
			&product.Barcode,
		)
		if err != nil {
			pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
			return nil, true
		}

		products = append(products, product)
		product.Image = []byte{}
		log.Println(product)
	}
	return products, false
}
