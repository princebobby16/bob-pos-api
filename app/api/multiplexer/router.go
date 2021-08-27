package multiplexer

import (
	"github.com/gorilla/mux"
	"gitlab.com/pbobby001/bobpos_api/app/api/handlers"
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
			Handler: handlers.HealthCheckHandler,
		},

		// products start
		Route{
			Name:    "Create ProductCreate",
			Path:    "/products",
			Method:  http.MethodPost,
			Handler: handlers.ProductCreate,
		},

		Route{
			Name:    "Delete ProductCreate",
			Path:    "/products",
			Method:  http.MethodDelete,
			Handler: handlers.DeleteProduct,
		},

		Route{
			Name:    "Get One ProductCreate By Id",
			Path:    "/one/products",
			Method:  http.MethodGet,
			Handler: handlers.GetOneProductById,
		},

		Route{
			Name:    "Get All Products",
			Path:    "/all/products",
			Method:  http.MethodGet,
			Handler: handlers.GetAllProducts,
		},

		Route{
			Name:    "Upload ProductCreate Image",
			Path:    "/up/products",
			Method:  http.MethodPost,
			Handler: handlers.HandleMediaUpload,
		},
		Route{
			Name:    "Delete Uploaded media file",
			Path:    "/can/products",
			Method:  http.MethodPost,
			Handler: handlers.HandleCancelMediaUpload,
		},
		// products end

		Route{
			Name:    "Get All Categories",
			Path:    "/categories",
			Method:  http.MethodGet,
			Handler: handlers.GetAllCategories,
		},

		Route{
			Name:    "Create Tax",
			Path:    "/tax",
			Method:  http.MethodGet,
			Handler: handlers.CreateTax,
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
