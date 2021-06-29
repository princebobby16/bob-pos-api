package pkg

import (
	"github.com/twinj/uuid"
	"time"
)

type (
	Product struct {
		Id            uuid.UUID       `json:"id"`
		Name          string          `json:"name"`
		Category      ProductCategory `json:"category"`
		Weight        string          `json:"weight"`
		CostPrice     float64         `json:"cost_price"`
		Tax           TaxDetails      `json:"tax"`
		ProfitMargin  float64         `json:"profit_margin"`
		Image         []byte          `json:"image"`
		NumberInStock int             `json:"number_in_stock"`
		Barcode       string          `json:"barcode"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
	}

	TaxDetails struct {
		Id        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		TaxRate   string    `json:"tax_rate"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	ProductCategory struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	GetAllCategoryResponse struct {
		Data []ProductCategory `json:"data"`
		Meta Meta              `json:"meta"`
	}

	StandardGetAllProductsResponse struct {
		Data []Product `json:"data"`
		Meta Meta      `json:"meta"`
	}

	GetOneProductResponse struct {
		Product Product `json:"product"`
		Meta    Meta    `json:"meta"`
	}

	StandardResponse struct {
		Data Data `json:"data"`
		Meta Meta `json:"meta"`
	}

	StandardCreatedProductResponse struct {
		Data CreatedProductData `json:"data"`
		Meta Meta               `json:"meta"`
	}

	CreatedProductData struct {
		Id             string `json:"id"`
		UiMessage      string `json:"ui_message"`
		ProductBarcode []byte `json:"product_barcode"`
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
