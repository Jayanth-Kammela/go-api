package routes

import (
	"net/http"

	"github.com/Jayanth-Kammela/go-api/controllers"
	"github.com/Jayanth-Kammela/go-api/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	//base
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	}).Methods("GET")

	// Products endpoints
	router.Use(middleware.ContentTypeMiddleware)
	router.HandleFunc("/api/v1/products", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/api/v1/product/{id}", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/api/v1/product", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/api/v1/product/{id}", controllers.UpdateProduct).Methods("PATCH")
	router.HandleFunc("/api/v1/product/{id}", controllers.DeleteProduct).Methods("DELETE")

	return router
}
