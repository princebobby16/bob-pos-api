package handlers

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/logger"
	"io/ioutil"
	"net/http"
)

func CreateTax(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	logger.Logger.Infof("Headers => TraceId: %s", traceId)

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
	}

	logger.Logger.Info("Request Object: ", string(requestBody))

	// Create ProductCreate instance to decode request object into
	var tax *pkg.TaxDetails

	// Decode request body into the Post struct
	err = json.Unmarshal(requestBody, &tax)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
	}
	logger.Logger.Info(tax)
}
