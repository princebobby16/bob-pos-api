package pkg

import (
	"encoding/json"
	"errors"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"time"
)

// ValidateHeadersAndReturnTheirValues Validate header is a function used to make sure that the required  headers are sent to the API
func ValidateHeadersAndReturnTheirValues(r *http.Request) (map[string]string, error) {
	//Group the headers
	receivedHeaders := make(map[string]string)
	requiredHeaders := []string{"trace-id"}

	for _, header := range requiredHeaders {
		value := r.Header.Get(header)
		if value != "" {
			receivedHeaders[header] = value
		} else if value == "" {
			return nil, errors.New("Required header: " + header + " not found")
		} else {
			return nil, errors.New("no headers received be sure to send some headers")
		}
	}

	return receivedHeaders, nil
}

// SendErrorResponse /* Helper func to handle error */
func SendErrorResponse(w http.ResponseWriter, tId uuid.UUID, traceId string, err error, httpStatus int) {
	w.WriteHeader(httpStatus)
	log.Println(err)
	_ = json.NewEncoder(w).Encode(StandardResponse{
		Data: Data{
			UiMessage: err.Error(),
		},
		Meta: Meta{
			Timestamp:     time.Now(),
			TransactionId: tId.String(),
			TraceId:       traceId,
			Status:        "FAILED",
		},
	})
	return
}
