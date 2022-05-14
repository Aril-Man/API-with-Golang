package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	ID        int             `json:"id" gorm:"primary_key"`
	Code      string          `json:"code" validate:"required"`
	Name      string          `json:"name" validate:"required"`
	Price     decimal.Decimal `json:"price" sql:"type:decimal(16,2)" validate:"required"`
	Deskripsi string          `json:"deskripsi" validate:"required"`
}

type Result struct {
	Code    int         `json:"code"`
	Count   int         `json:"count"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/api-go?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("Connected!")
	}

	db.AutoMigrate(&Product{})

	handleRequests()
}

func handleRequests() {
	log.Println("Mulai Development Server di http://127.0.0.1:9999")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)
	router.HandleFunc("/products", getProducts).Methods("GET")
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payload, &product)

	validate := validator.New()

	if err := validate.Struct(product); err != nil {
		res := Result{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data:    err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res.Data.(string)))
	} else {

		db.Create(&product)

		res := Result{
			Code:    http.StatusCreated,
			Data:    product,
			Message: "Product berhasil dibuat",
		}

		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}

}

func getProducts(w http.ResponseWriter, r *http.Request) {
	product := []Product{}

	db.Find(&product)

	if len(product) <= 0 {
		res := Result{
			Code:    http.StatusNotFound,
			Message: "Product Kosong"}

		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	} else {

		res := Result{
			Code:    http.StatusOK,
			Count:   len(product),
			Data:    product,
			Message: "Product berhasil ditemukan"}

		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}

}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	db.First(&product, id)

	if product.ID == 0 {
		res := Result{
			Code:    http.StatusNotFound,
			Message: "Product tidak ditemukan"}

		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	} else {

		res := Result{
			Code:    http.StatusOK,
			Data:    product,
			Message: "Product berhasil ditemukan"}

		result, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)
	id := vals["id"]

	var product Product
	db.First(&product, id)

	db.Delete(&product)

	res := Result{
		Code:    http.StatusOK,
		Data:    product,
		Message: "Product berhasil dihapus"}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)
	id := vals["id"]

	payload, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payload, &product)

	db.Model(&product).Where("id = ?", id).Update(&product)

	res := Result{
		Code:    http.StatusOK,
		Data:    product,
		Message: "Product berhasil diupdate"}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}
