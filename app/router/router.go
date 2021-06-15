package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/pbobby001/bobpos_api/app/controllers"
	"gitlab.com/pbobby001/bobpos_api/app/controllers/mediaupload"
	"gitlab.com/pbobby001/bobpos_api/app/controllers/product"
	"net/http"
)

//Route Create a single route object
type Route struct {
	Name    string
	Path    string
	Method  string
	Handler http.HandlerFunc
}

//Routes Create an object of different routes
type Routes []Route

// InitRoutes Set up routes
func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	routes := Routes{
		// health check
		Route{
			Name:    "Health Check",
			Path:    "/",
			Method:  http.MethodGet,
			Handler: controllers.HealthCheckHandler,
		},

		Route{
			Name:    "Create Product",
			Path:    "/products",
			Method:  http.MethodPost,
			Handler: product.CreateProduct,
		},

		Route{
			Name:    "Delete Product",
			Path:    "/products",
			Method:  http.MethodDelete,
			Handler: product.DeleteProduct,
		},

		Route{
			Name:    "Get All Products",
			Path:    "/products",
			Method:  http.MethodGet,
			Handler: product.GetAllProducts,
		},

		Route{
			Name:    "Upload Product Image",
			Path:    "/up/products",
			Method:  http.MethodPost,
			Handler: mediaupload.HandleMediaUpload,
		},
		Route{
			Name:    "Delete Uploaded media file",
			Path:    "/can/products",
			Method:  http.MethodPost,
			Handler: mediaupload.HandleCancelMediaUpload,
		},
	}

	for _, route := range routes {
		router.Name(route.Name).
			Methods(route.Method).
			Path(route.Path).
			Handler(route.Handler)
	}

	return router
}
