package main

import (
	"ecommerce/middleware"
	"ecommerce/products"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	err := products.NewProductDB()
	if err != nil {
		fmt.Println("Failed to create ")
		os.Exit(1)
	}

	r.HandleFunc("/auth", middleware.Authentication).Methods("POST")
	///The r.Use() function takes a middleware function as an argument. This middleware function must match the mux.MiddlewareFunc type,
	//which is defined as a function that takes an http.Handler and returns an http.Handler.
	// Essentially, it wraps the original handler with additional functionality.
	//r.Use(AuthMiddleWare)

	//HandleFunc registers a handler function for a given pattern in the URL path.
	//The handler function must have the signature func(w http.ResponseWriter, r *http.Request), where w is used to write the response and r contains the request details.
	//This is a convenient way to quickly set up routes without needing to create a specific type or object that implements the http.Handler interface.`

	//Handle, on the other hand, requires an object that implements the http.Handler interface. This interface defines a single method,
	//ServeHTTP, which has the same signature as the function used with HandleFunc. Using Handle allows for more structured and object-oriented handling of HTTP requests,
	// which can be useful for complex applications where you might want to encapsulate request handling logic within different types.
	r.Handle("/product", middleware.AuthMiddleWare(http.HandlerFunc(products.CreateProduct))).Methods("POST")
	r.Handle("/product", middleware.AuthMiddleWare(http.HandlerFunc(products.UpdateProduct))).Methods("PUT")
	r.Handle("/product/{id}", middleware.AuthMiddleWare(http.HandlerFunc(products.DeleteProduct))).Methods("DELETE")
	r.Handle("/product/{id}", middleware.AuthMiddleWare(http.HandlerFunc(products.GetProduct))).Methods("GET")
	r.Handle("/product", middleware.AuthMiddleWare(http.HandlerFunc(products.GetProducts))).Methods("GET")
	fmt.Println("Listening on 9000")
	http.ListenAndServe(":9000", r)

}
