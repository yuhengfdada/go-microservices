package data

import (
	"encoding/json"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// Product ID
	// required: true
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", isValidSKU)
	err := validate.Struct(p)
	return err
}

func isValidSKU(fl validator.FieldLevel) bool {
	val := fl.Field().Interface().(string)
	reg := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]+")
	matches := reg.FindAllString(val, -1)
	return len(matches) == 1
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

func AddProduct(prod *Product) {
	prod.ID = nextID()
	productList = append(productList, prod)
}

func UpdateProduct(id int, prod *Product) {
	prod.ID = id
	productList[id-1] = prod
}

func DeleteProduct(id int) {
	var idx int = -1
	for i, prod := range productList {
		if prod.ID == id {
			idx = i
		}
	}
	if idx == -1 {
		return
	}
	productList = append(productList[:idx], productList[idx+1:]...)
}

func nextID() int {
	if len(productList) == 0 {
		return 1
	}
	return productList[len(productList)-1].ID + 1
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
