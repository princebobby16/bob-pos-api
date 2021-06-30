package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/pbobby001/bobpos_api/app/controllers"
	"gitlab.com/pbobby001/bobpos_api/app/controllers/category"
	"gitlab.com/pbobby001/bobpos_api/app/controllers/mediaupload"
	"gitlab.com/pbobby001/bobpos_api/app/controllers/product"
	"gitlab.com/pbobby001/bobpos_api/app/controllers/product/create"
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

		// products start
		Route{
			Name:    "Create ProductCreate",
			Path:    "/products",
			Method:  http.MethodPost,
			Handler: create.ProductCreate,
		},

		Route{
			Name:    "Delete ProductCreate",
			Path:    "/products",
			Method:  http.MethodDelete,
			Handler: product.DeleteProduct,
		},

		Route{
			Name:    "Get One ProductCreate By Id",
			Path:    "/one/products",
			Method:  http.MethodGet,
			Handler: product.GetOneProductById,
		},

		Route{
			Name:    "Get All Products",
			Path:    "/all/products",
			Method:  http.MethodGet,
			Handler: product.GetAllProducts,
		},

		Route{
			Name:    "Upload ProductCreate Image",
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
		// products end

		Route{
			Name:    "Get All Categories",
			Path:    "/categories",
			Method:  http.MethodGet,
			Handler: category.GetAllCategories,
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
