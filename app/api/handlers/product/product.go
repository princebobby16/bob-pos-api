package product

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"log"
	"net/http"
	"time"
)

func ProductCreate(w http.ResponseWriter, r *http.Request) {
	transactionId, err, traceId, product, done := handleCreatProductRequest(w, r)
	if done {
		return
	}

	if getImageIfAvailable(w, transactionId, traceId, err, product) {
		return
	}

	productId, done := insertProductIntoDatabase(w, err, product, transactionId, traceId)
	if done {
		return
	}

	sendProductCreatedSuccessResponse(w, productId, transactionId, traceId)
}

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
	log.Printf("Headers => TraceId: %s", traceId)

	productId := r.URL.Query().Get("product_id")
	log.Println(productId)

	product, done := getProductFromDatabase(w, err, productId, transactionId)
	if done {
		return
	}

	sendSuccessResponse(w, product, transactionId, traceId)
}

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

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	transactionId, err, traceId, done := generateTransactionIdAndExtractTraceId(w, r)
	if done {
		return
	}

	categories, done2 := getAllCategoriesFromTheDatabase(w, err, transactionId, traceId)
	if done2 {
		return
	}
	log.Println(categories)

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
