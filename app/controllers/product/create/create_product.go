package create

import (
	"net/http"
)

func ProductCreate(w http.ResponseWriter, r *http.Request) {
	transactionId, err, traceId, product, done := handleCreatProductRequest(w, r)
	if done {
		return
	}

	if getImageIfAvailable(w, transactionId, traceId, err, product) {
		return
	}

	productId, done3 := insertProductIntoDatabase(w, err, product, transactionId, traceId)
	if done3 {
		return
	}

	sendProductCreatedSuccessResponse(w, productId, transactionId, traceId)
}
