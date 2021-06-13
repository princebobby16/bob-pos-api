package pkg

import (
	"encoding/json"
	"errors"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg/logs"
	"net/http"
	"time"
)

// ValidateHeaders /* Validate header is a function used to make sure that the required  headers are sent to the API
func ValidateHeaders(r *http.Request) (map[string]string, error) {
	//Group the headers
	receivedHeaders := make(map[string]string)
	requiredHeaders := []string{"trace-id", "tenant-namespace"}

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

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

// SendErrorResponse /* Helper func to handle error */
func SendErrorResponse(w http.ResponseWriter, tId uuid.UUID, traceId string, err error, httpStatus int) {
	w.WriteHeader(httpStatus)
	_ = logs.Logger.Error(err)
	_ = json.NewEncoder(w).Encode(StandardResponse{
		Data: Data{
			Id:        "",
			UiMessage: "Something went wrong! Contact Admin!",
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