package product

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/db"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/logs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	transactionId, err, traceId, product, done := handleCreatProductRequest(w, r)
	if done {
		return
	}

	if getImageIfAvailable(err, product) {
		return
	}

	query := `insert into bobpos.products(id, name, category, weight, cost_price, tax, profit_margin, image, number_in_stock)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	productId := uuid.NewV4()
	_, err = db.Connection.Exec(
		query,
		&productId,
		&product.Name,
		&product.Category,
		&product.Weight,
		&product.CostPrice,
		&product.Tax,
		&product.ProfitMargin,
		&product.Image,
		&product.NumberInStock,
	)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        productId.String(),
			UiMessage: "Product Created!",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})

}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeaders(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	logs.Logger.Infof("Headers => TraceId: %s", traceId)

	productId := r.URL.Query().Get("product_id")
	logs.Logger.Info(productId)

	query := `delete from bobpos.products where id=$1`

	_, err = db.Connection.Exec(query, &productId)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        productId,
			UiMessage: "Product Deleted!",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})

}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	logs.Logger.Info("in get")
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeaders(r)
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

func GetOneProductById(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeaders(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	logs.Logger.Infof("Headers => TraceId: %s", traceId)

	productId := r.URL.Query().Get("product_id")
	logs.Logger.Info(productId)

	query := `select * from bobpos.products where id=$1`
	var product pkg.Product
	err = db.Connection.QueryRow(query, &productId).Scan(
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
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pkg.GetOneProductResponse{
		Product: product,
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func getImageIfAvailable(err error, product *pkg.Product) bool {
	wd, err := os.Getwd()
	if err != nil {
		_ = logs.Logger.Error(err)
		return true
	}

	path := filepath.Join(wd, "pkg/images")
	logs.Logger.Info(path)

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			_ = logs.Logger.Warn(err)
		} else {
			_ = logs.Logger.Error(err)
			return true
		}
	}

	var imageBytes []byte

	if fileInfo != nil {
		for _, file := range fileInfo {
			logs.Logger.Info(file.Name())

			fileLocation := filepath.Join(path, file.Name())

			openImage, err := os.Open(fileLocation)
			if err != nil {
				_ = logs.Logger.Error(err)
				return true
			}

			imageBytes, err = ioutil.ReadAll(openImage)
			if err != nil {
				_ = logs.Logger.Error(err)
				return true
			}
			err = openImage.Close()
			if err != nil {
				_ = logs.Logger.Error(err)
				return true
			}
			product.Image = imageBytes

			err = os.RemoveAll(wd + "/pkg/images")
			if err != nil {
				_ = logs.Logger.Error(err)
				return true
			}
		}
	}
	return false
}

func handleCreatProductRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, error, string, *pkg.Product, bool) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeaders(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	logs.Logger.Infof("Headers => TraceId: %s", traceId)

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}

	logs.Logger.Info("Request Object: ", string(requestBody))

	// Create Product instance to decode request object into
	var product *pkg.Product

	// Decode request body into the Post struct
	err = json.Unmarshal(requestBody, &product)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}
	logs.Logger.Info(product)
	return transactionId, err, traceId, product, false
}
