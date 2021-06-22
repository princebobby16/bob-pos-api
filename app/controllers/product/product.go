package product

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/db"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"gitlab.com/pbobby001/bobpos_api/pkg/logs"
	"image/png"
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
	// create barcode
	imageData, code, done2 := generateBarcodeForProduct(w, err, transactionId, traceId)
	if done2 {
		return
	}

	if getImageIfAvailable(w, transactionId, traceId, err, product) {
		return
	}
	productId, done3 := insertProductIntoDatabase(w, err, code, product, transactionId, traceId)
	if done3 {
		return
	}

	sendProductCreatedSuccessResponse(w, imageData, productId, transactionId, traceId)
}

func sendProductCreatedSuccessResponse(w http.ResponseWriter, imageData []byte, productId uuid.UUID, transactionId uuid.UUID, traceId string) {
	_ = json.NewEncoder(w).Encode(pkg.StandardCreatedProductResponse{
		Data: pkg.CreatedProductData{
			Id:             productId.String(),
			UiMessage:      "Product Created!",
			ProductBarcode: imageData,
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func insertProductIntoDatabase(w http.ResponseWriter, err error, code string, product *pkg.Product, transactionId uuid.UUID, traceId string) (uuid.UUID, bool) {
	query := `insert into bobpos.products(id, name, category, weight, cost_price, tax, profit_margin, image, number_in_stock, barcode)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

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
		&code,
	)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return uuid.UUID{}, true
	}
	return productId, false
}

func generateRandomNumbers(num int, otpChars string) (string, error) {
	buffer := make([]byte, num)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < num; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func generateBarcodeForProduct(w http.ResponseWriter, err error, transactionId uuid.UUID, traceId string) ([]byte, string, bool) {
	barcodeNumbers, bytes, s, b, done := getBarcodeNumber(w, err, transactionId, traceId)
	if done {
		return bytes, s, b
	}

	content, file, i, s2, b2, done2 := createBarcodeImage(w, err, transactionId, traceId, barcodeNumbers)
	if done2 {
		return i, s2, b2
	}
	imageBytes, i2, s3, b3, done3 := openAndReadBarcodeContent(w, err, transactionId, traceId, file)
	if done3 {
		return i2, s3, b3
	}
	return imageBytes, content, false
}

func openAndReadBarcodeContent(w http.ResponseWriter, err error, transactionId uuid.UUID, traceId string, file *os.File) ([]byte, []byte, string, bool, bool) {
	f, err := os.Open(file.Name())
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return nil, nil, "", true, true
	}

	imageBytes, err := ioutil.ReadAll(f)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return nil, nil, "", true, true
	}
	err = file.Close()
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return nil, nil, "", true, true
	}
	return imageBytes, nil, "", false, false
}

func createBarcodeImage(w http.ResponseWriter, err error, transactionId uuid.UUID, traceId string, barcodeNumbers string) (string, *os.File, []byte, string, bool, bool) {
	barCode, err := ean.Encode(barcodeNumbers)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return "", nil, nil, "", true, true
	}

	content := barCode.Content()
	logs.Logger.Info(content)

	code, err := barcode.Scale(barCode, 200, 100)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return "", nil, nil, "", true, true
	}

	file, err := os.Create("qrcode.png")
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return "", nil, nil, "", true, true
	}

	//encode the barcode as png
	_ = png.Encode(file, code)
	return content, file, nil, "", false, false
}

func getBarcodeNumber(w http.ResponseWriter, err error, transactionId uuid.UUID, traceId string) (string, []byte, string, bool, bool) {
	last2Numbers, err := generateRandomNumbers(2, "1234567890")
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return "", nil, "", true, true
	}
	middleNumbers, err := generateRandomNumbers(5, "567890")
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, traceId, err, http.StatusInternalServerError)
		return "", nil, "", true, true
	}
	barcodeNumbers := "00233" + middleNumbers + last2Numbers
	return barcodeNumbers, nil, "", false, false
}

func getImageIfAvailable(w http.ResponseWriter, tid uuid.UUID, traceId string, err error, product *pkg.Product) bool {
	wd, err := os.Getwd()
	if err != nil {
		_ = logs.Logger.Error(err)
		pkg.SendErrorResponse(w, tid, traceId, err, http.StatusBadRequest)
		return true
	}

	path := filepath.Join(wd, "pkg/images")
	logs.Logger.Info(path)

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			_ = logs.Logger.Warn(err)
			pkg.SendErrorResponse(w, tid, traceId, errors.New("no image uploaded"), http.StatusBadRequest)
			return true
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
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}

			imageBytes, err = ioutil.ReadAll(openImage)
			if err != nil {
				_ = logs.Logger.Error(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}
			err = openImage.Close()
			if err != nil {
				_ = logs.Logger.Error(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}
			product.Image = imageBytes

			err = os.RemoveAll(wd + "/pkg/images")
			if err != nil {
				_ = logs.Logger.Error(err)
				pkg.SendErrorResponse(w, tid, traceId, err, http.StatusInternalServerError)
				return true
			}
			err = os.Remove(wd + "/qrcode.png")
			if err != nil {
				_ = logs.Logger.Error(err)
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
