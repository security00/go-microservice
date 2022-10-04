package data

import "testing"

func TestProductValidate(t *testing.T) {
	p := &Product{
		Name:  "Water",
		Price: 10,
		SKU:   "a-b",
	}

	err := p.Validate()

	if err != nil {
		t.Log(err)
	}
}
