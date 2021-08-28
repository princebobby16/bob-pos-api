package handlers

import (
	"encoding/json"
	"errors"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/db/connection"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func sendProductCreatedSuccessResponse(w http.ResponseWriter, productId uuid.UUID, transactionId uuid.UUID, traceId string) {
	_ = json.NewEncoder(w).Encode(pkg.StandardCreatedProductResponse{
		Data: pkg.CreatedProductData{
			Id:        productId.String(),
			UiMessage: "ProductCreate Created!",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func insertProductIntoDatabase(w http.ResponseWriter, err error, product *pkg.Product, transactionId uuid.UUID, traceId string) (uuid.UUID, bool) {

	query := `select count(name) from bobpos.tax where name = $1;`
	var count int
	err1 := connection.Connection.QueryRow(query, product.Tax.Name).Scan(&count)
	if err1 != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return uuid.UUID{}, true
	}

	if count == 0 {
		pkg.SendErrorResponse(w, transactionId, traceId, errors.New("tax does not exist"), http.StatusInternalServerError)
		return uuid.UUID{}, true
	}

	query = `insert into bobpos.products(id, name, category, weight, cost_price, tax, profit_margin, image, number_in_stock, barcode)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	productId := uuid.NewV4()
	_, err = connection.Connection.Exec(
		query,
		&productId,
		&product.Name,
		&product.Category.Name,
		&product.Weight,
		&product.CostPrice,
		&product.Tax.TaxRate,
		&product.ProfitMargin,
		&product.Image,
		&product.NumberInStock,
		&product.Barcode,
	)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return uuid.UUID{}, true
	}
	return productId, false
}

func getImageIfAvailable(w http.ResponseWriter, tid uuid.UUID, traceId string, err error, product *pkg.Product) bool {
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		pkg.SendErrorResponse(w, tid, traceId, err, http.StatusBadRequest)
		return true
	}

	path := filepath.Join(wd, "pkg/images")
	log.Println(path)

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(err)
			pkg.SendErrorResponse(w, tid, traceId, errors.New("no image uploaded"), http.StatusBadRequest)
			return true
		} else {
			log.Println(err)
			return true
		}
	}

	var imageBytes []byte

	if fileInfo != nil {
		for _, file := range fileInfo {
			log.Println(file.Name())

			fileLocation := filepath.Join(path, file.Name())

			openImage, err := os.Open(fileLocation)
			if err != nil {
				log.Println(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}

			imageBytes, err = ioutil.ReadAll(openImage)
			if err != nil {
				log.Println(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}
			err = openImage.Close()
			if err != nil {
				log.Println(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}
			product.Image = imageBytes

			err = os.RemoveAll(wd + "/pkg/images")
			if err != nil {
				log.Println(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}
		}
	}
	return false
}

func handleCreatProductRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, error, string, *pkg.Product, bool) {

	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	log.Printf("Headers => TraceId: %s", traceId)

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}

	log.Println("Request Object: ", string(requestBody))

	// Create ProductCreate instance to decode request object into
	var product *pkg.Product

	// Decode request body into the Post struct
	err = json.Unmarshal(requestBody, &product)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}
	log.Println(product)

	if product.Barcode == "" {
		pkg.SendErrorResponse(w, transactionId, traceId, errors.New("no barcode provided for product"), http.StatusBadRequest)
		return uuid.UUID{}, nil, "", nil, true
	}

	return transactionId, err, traceId, product, false
}
