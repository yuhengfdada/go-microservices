package data

import (
	"testing"
)

func TestValidate(t *testing.T) {
	prod := &Product{
		ID:    3,
		Name:  "tea",
		Price: 0.1,
		SKU:   "adas-aa-aa",
	}
	err := prod.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
