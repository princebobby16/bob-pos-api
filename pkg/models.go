package pkg

import (
	"github.com/twinj/uuid"
	"time"
)

type (
	Product struct {
		Id            uuid.UUID `json:"id"`
		Name          string    `json:"name"`
		Category      string    `json:"category"`
		Weight        string    `json:"weight"`
		Price         int       `json:"price"`
		Image         []byte    `json:"image"`
		NumberInStock int       `json:"number_in_stock"`
	}

	StandardResponse struct {
		Data Data `json:"data"`
		Meta Meta `json:"meta"`
	}

	Data struct {
		Id        string `json:"id"`
		UiMessage string `json:"ui_message"`
	}

	Meta struct {
		Timestamp     time.Time `json:"timestamp"`
		TransactionId string    `json:"transaction_id"`
		TraceId       string    `json:"trace_id"`
		Status        string    `json:"status"`
	}
)
