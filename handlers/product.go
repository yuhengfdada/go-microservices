package handlers

import (
	"log"
	"net/http"

	"github.com/yuhengfdada/go-microservices/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Product called")
	if r.Method == http.MethodGet {
		p.getProducts(rw)
		return
	}
	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}
	http.Error(rw, "Method not supported", http.StatusMethodNotAllowed)
}

func (p *Product) getProducts(rw http.ResponseWriter) {
	p.l.Println("GetProducts() called")
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		p.l.Fatalln(err)
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (p *Product) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("AddProducts() called")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Added product: %#v", prod)
}
