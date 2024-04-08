package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"ecommerce/database"
	"ecommerce/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ProductRepo struct {
	Db *gorm.DB
}

func NewProductDB() error {
	db, err := database.InitDb()
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.Product{})
	return err
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println("Error decoding the json ", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = models.CreateProduct(database.DB, &product)
	if err != nil {
		fmt.Println("Error creating the product", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, product)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	var product_id int
	var err error
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok {
		product_id, err = strconv.Atoi(id)
		if err != nil {
			fmt.Println("Unable to convert id to int", err)
			respondWithError(w, http.StatusBadRequest, "Invalid product id")
			return
		}
	}
	err = models.GetProduct(database.DB, &product, product_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondWithError(w, http.StatusBadRequest, "Record not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println("Error decoding the json ", err)

	}
	err = models.UpdateProduct(database.DB, &product)
	if err != nil {
		fmt.Println("Error creating the product", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())

	}
	respondWithJSON(w, http.StatusOK, product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	var product_id int
	var err error
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok {
		product_id, err = strconv.Atoi(id)
		if err != nil {
			fmt.Println("Unable to convert id to int", err)
			respondWithError(w, http.StatusBadRequest, "Invalid product id")
		}
	}
	err = models.DeleteProduct(database.DB, &product, product_id)
	if err != nil {
		fmt.Println("Error deleting the product", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("##### Inside get product #####")
	var product models.Product
	var product_id int
	var err error
	vars := mux.Vars(r)
	id, ok := vars["id"]
	fmt.Println("product id", id)
	if ok {
		product_id, err = strconv.Atoi(id)
		if err != nil {
			fmt.Println("Unable to convert id to int", err)
			respondWithError(w, http.StatusBadRequest, "Invalid product id")
		}
	}
	err = models.GetProduct(database.DB, &product, product_id)
	if err != nil {
		fmt.Println("Error getting the product", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	err := models.GetProducts(database.DB, &products)
	if err != nil {
		fmt.Println("Error getting products list")
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, products)
}
