package media

import (
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/twinj/uuid"
	"gitlab.com/pbobby001/bobpos_api/pkg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	_ int = iota
	_     = 1 << (10 * iota)
	MB
)

func HandleMediaUpload(w http.ResponseWriter, r *http.Request) {

	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	log.Println("Headers => TraceId: " + traceId)

	err = r.ParseMultipartForm(10 * MB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	file, handler, err := r.FormFile("media_file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	log.Println("Uploaded File: ", handler.Filename)
	log.Println("File Size: ", handler.Size)
	log.Println("MIME Header: ", handler.Header)

	if strings.Split(handler.Filename, ".")[1] == "webp" {
		log.Println("invalid image format")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid\timage\tformat"))
		return
	}

	f := make(chan multipart.File)
	go parseMultipartToFile(f, handler.Filename)
	f <- file
	close(f)

	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        "",
			UiMessage: "FILE UPLOADING",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func parseMultipartToFile(fileChannel <-chan multipart.File, filename string) {
	// Listen on the file channel
	for file := range fileChannel {

		// read the file bytes
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			return
		}

		// get the working directory for generating the path for image storage
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}

		// join the working directory path with the path for image storage
		join := filepath.Join(wd, "pkg/images")

		// create a new directory for storing the image
		err = os.Mkdir(join, 0755)
		if err != nil {
			if os.IsExist(err) {
				log.Println(err)
			} else {
				log.Println(err)
				return
			}
		}

		tempFile, err := os.Create(join + "/" + filename)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = tempFile.Write(fileBytes)
		if err != nil {
			log.Println(err)
			return
		}

		img, err := imaging.Open(tempFile.Name())
		if err != nil {
			log.Println(err)
			return
		}

		imb := imaging.AdjustBrightness(img, -5)
		src := imaging.Resize(imb, 0, 200, imaging.Lanczos)
		err = imaging.Save(src, tempFile.Name())
		if err != nil {
			log.Println(err)
			return
		}
		_ = tempFile.Close()

		log.Println("Successfully resized image...")
	}
}

func HandleCancelMediaUpload(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	//Get the relevant headers
	traceId := headers["trace-id"]

	// Logging the headers
	log.Println("Headers => TraceId: " + traceId)

	fileName := r.URL.Query().Get("file_name")
	log.Println(fileName)

	workingDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(workingDir)

	imageStoragePath := path.Join(workingDir, "/pkg/images")
	log.Println(imageStoragePath)

	if imageStoragePath == "" {
		log.Println("No image has been uploaded to server")
		w.WriteHeader(http.StatusGone)
		_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
			Data: pkg.Data{
				Id:        "",
				UiMessage: "no file found, please upload image!",
			},
			Meta: pkg.Meta{
				Timestamp:     time.Now(),
				TransactionId: transactionId.String(),
				TraceId:       traceId,
				Status:        "SUCCESS",
			},
		})
		return
	}

	err = os.Remove(imageStoragePath + "/" + fileName)
	if err != nil {
		log.Println(err)
		return
	}

	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        "",
			UiMessage: "FILE DELETED",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}

func DeleteUploadedFiles(w http.ResponseWriter, r *http.Request) {
	transactionId := uuid.NewV4()

	headers, err := pkg.ValidateHeadersAndReturnTheirValues(r)
	if err != nil {
		pkg.SendErrorResponse(w, transactionId, "", err, http.StatusBadRequest)
		return
	}

	//Get the relevant headers
	traceId := headers["trace-id"]
	tenantNamespace := headers["tenant-namespace"]

	// Logging the headers
	log.Println("Headers => TraceId: " + traceId + ", TenantNamespace: " + tenantNamespace)

	workingDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(workingDir)

	imageStoragePath := path.Join(workingDir, "/pkg/images")
	log.Println(imageStoragePath)

	if imageStoragePath == "" {
		log.Println("No image has been uploaded to server")
		w.WriteHeader(http.StatusGone)
		_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
			Data: pkg.Data{
				Id:        "",
				UiMessage: "no file found, please upload image!",
			},
			Meta: pkg.Meta{
				Timestamp:     time.Now(),
				TransactionId: transactionId.String(),
				TraceId:       traceId,
				Status:        "SUCCESS",
			},
		})
		return
	}

	err = os.RemoveAll(imageStoragePath)
	if err != nil {
		log.Println(err)
		if os.IsNotExist(err) {

		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
				Data: pkg.Data{
					Id:        "",
					UiMessage: "something went wrong! contact admin",
				},
				Meta: pkg.Meta{
					Timestamp:     time.Now(),
					TransactionId: transactionId.String(),
					TraceId:       traceId,
					Status:        "SUCCESS",
				},
			})
			return
		}
	}

	_ = json.NewEncoder(w).Encode(pkg.StandardResponse{
		Data: pkg.Data{
			Id:        "",
			UiMessage: "FILES DELETED",
		},
		Meta: pkg.Meta{
			Timestamp:     time.Now(),
			TransactionId: transactionId.String(),
			TraceId:       traceId,
			Status:        "SUCCESS",
		},
	})
}
