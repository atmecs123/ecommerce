package products

import (
	"ecommerce/database"
	"ecommerce/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var productLog logr.Logger

type ProductRepo struct {
	Db *gorm.DB
}

func NewProductDB() error {
	productLog = productLog.WithValues("product")
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
	ctx := r.Context()
	productLog.Info("Enter create product")
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		productLog.Error(err, "Error decoding the json ")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = models.CreateProduct(ctx, database.DB, &product)
	if err != nil {
		productLog.Error(err, "Error creating the product")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, product)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productLog.Info("Enter update product")
	ctx := r.Context()
	var product models.Product
	var product_id int
	var err error
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok {
		product_id, err = strconv.Atoi(id)
		if err != nil {
			productLog.Error(err, "Unable to convert id to int")
			respondWithError(w, http.StatusBadRequest, "Invalid product id")
			return
		}
	}
	err = models.GetProduct(ctx, database.DB, &product, product_id)
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
		productLog.Error(err, "Error decoding the json ")

	}
	err = models.UpdateProduct(ctx, database.DB, &product)
	if err != nil {
		productLog.Error(err, "Error creating the product")
		respondWithError(w, http.StatusInternalServerError, err.Error())

	}
	respondWithJSON(w, http.StatusOK, product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productLog.Info("Enter delete product")
	ctx := r.Context()
	var product models.Product
	var product_id int
	var err error
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok {
		product_id, err = strconv.Atoi(id)
		if err != nil {
			productLog.Error(err, "Unable to convert id to int")
			respondWithError(w, http.StatusBadRequest, "Invalid product id")
		}
	}
	err = models.DeleteProduct(ctx, database.DB, &product, product_id)
	if err != nil {
		productLog.Error(err, "Error deleting the product")
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productLog.Info("Enter get product")
	ctx := r.Context()
	var product models.Product
	var product_id int
	var err error
	vars := mux.Vars(r)
	id, ok := vars["id"]
	fmt.Println("product id", id)
	if ok {
		product_id, err = strconv.Atoi(id)
		if err != nil {
			productLog.Error(err, "Unable to convert id to int", err)
			respondWithError(w, http.StatusBadRequest, "Invalid product id")
		}
	}
	err = models.GetProduct(ctx, database.DB, &product, product_id)
	if err != nil {
		productLog.Error(err, "Error getting the product", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	productLog.Info("Enter get products")
	var products []models.Product
	ctx := r.Context()
	defaultpage := 1
	defaultlimit := 5
	query := r.URL.Query()
	pageParam := query.Get("page")
	limitParam := query.Get("limit")
	err := models.GetProducts(ctx, database.DB, &products, 0, 0)
	if err != nil {
		productLog.Error(err, "Error getting products list")
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	if pageParam == "" && limitParam == "" {
		//No pagination parameters return all elements
		respondWithJSON(w, http.StatusOK, products)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid limit number")
		return
	}

	if page < 1 {
		page = defaultpage
	}
	if limit < 1 {
		limit = defaultlimit
	}
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if startIndex > len(products) {
		startIndex = len(products)
	}
	if endIndex > len(products) {
		endIndex = len(products)
	}
	fmt.Println("start index", startIndex)
	fmt.Println("start index", endIndex)
	fmt.Println("products", products)
	paginatedProducts := products[startIndex:endIndex]
	respondWithJSON(w, http.StatusOK, paginatedProducts)
}

func GetProductsPaginated(w http.ResponseWriter, r *http.Request) {
	productLog.Info("Enter get products")
	ctx := r.Context()
	var products []models.Product
	//We extract page and limit parameters from the query string. These determine which portion of the data to send back.
	//The logic then calculates the start and end indexes based on these parameters and slices the items accordingly.
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	err := models.GetProducts(ctx, database.DB, &products, offset, limit)
	if err != nil {
		productLog.Error(err, "Error getting products list")
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	fmt.Println("##### context value is #######", ctx.Value("uuid"))
	respondWithJSON(w, http.StatusOK, products)
}
