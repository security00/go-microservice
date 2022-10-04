package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/security00/go-microservice/product-api/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{
		l: l,
	}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProduct()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := p.handleData(rw, r)

	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Convert to int failed", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT Product ", id)

	prod := p.handleData(rw, r)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Convert to int failed", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle DELETE Product ", id)

	err = data.DeleteProduct(id)
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

}

func (p *Products) handleData(rw http.ResponseWriter, r *http.Request) data.Product {
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	return prod
}

type KeyProduct struct {
}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		log.Printf("r.Body: %#v\n", r.Body)
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[Error] deserializing product ", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Println("[Error] validator product ", err)
			http.Error(rw, fmt.Sprintf("Validate product error: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
