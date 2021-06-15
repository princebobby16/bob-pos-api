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
		CostPrice     float64   `json:"cost_price"`
		Tax           float64   `json:"tax"`
		ProfitMargin  float64   `json:"profit_margin"`
		Image         []byte    `json:"image"`
		NumberInStock int       `json:"number_in_stock"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	StandardGetAllProductsResponse struct {
		Data []Product `json:"data"`
		Meta Meta      `json:"meta"`
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
