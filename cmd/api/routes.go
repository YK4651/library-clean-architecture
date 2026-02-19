package main

import (
	"github.com/YK4651/library-clean-architecture/internal/http/controllers"
	"github.com/gorilla/mux"
)

func setupRoutes(r *mux.Router, controller *controllers.BookController) {
	// パスパラメータを持つルート
	r.HandleFunc("/books", controller.ListBooks).Methods("GET")
	r.HandleFunc("/books/{bookID}", controller.GetBook).Methods("GET")

	// APIプレフィックス
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/books", controller.ListBooks).Methods("GET")
	api.HandleFunc("/books/{bookID}", controller.GetBook).Methods("GET")
}
