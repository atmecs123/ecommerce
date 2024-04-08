package main

import (
	"fmt"
	"ecommerce/products"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	///The application will be a simple REST API server
	///that will expose endpoints to allow accessing and manipulating 'products'. The operations that our endpoint will allow include:
	r := mux.NewRouter()
	err := products.NewProductDB()
	if err != nil {
		fmt.Println("Failed to create ")
	}
	r.HandleFunc("/product", products.CreateProduct).Methods("POST")
	r.HandleFunc("/product", products.UpdateProduct).Methods("PUT")
	r.HandleFunc("/product/{id}", products.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/product/{id}", products.GetProduct).Methods("GET")
	r.HandleFunc("/product", products.GetProducts).Methods("GET")
	fmt.Println("Listening on 9000")
	http.ListenAndServe(":9000", r)

}
