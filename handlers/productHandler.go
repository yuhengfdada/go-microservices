package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yuhengfdada/go-microservices/data"
)

type ProductHandler struct {
	l *log.Logger
}

func NewProductHandler(l *log.Logger) *ProductHandler {
	return &ProductHandler{l}
}

func (p *ProductHandler) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("GetProducts() called")
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		p.l.Fatalln(err)
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (p *ProductHandler) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AddProducts() called")
	prod := r.Context().Value(ProdKey{}).(*data.Product)
	data.AddProduct(prod)
	p.l.Printf("Added product: %#v", prod)
}

func (p *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("UpdateProducts() called")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Error converting id", http.StatusBadRequest)
	}
	prod := r.Context().Value(ProdKey{}).(*data.Product)
	data.UpdateProduct(id, prod)
}

type ProdKey struct{} // a key for retrieving the product in http.request.Context().

func (p *ProductHandler) MiddlewareProductConversion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{} // prod is type *data.Product
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Bad request: error unmarshalling json", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Bad request: error validating product, %s", err), http.StatusBadRequest)
			return
		}
		context := context.WithValue(r.Context(), ProdKey{}, prod)
		r = r.WithContext(context)
		next.ServeHTTP(rw, r)
	})
}
