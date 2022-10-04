package handlers

import (
	"github.com/security00/go-microservice/data"
	"log"
	"net/http"
	"regexp"
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

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		id := p.getId(rw, r)
		p.updateProduct(rw, r, id)
		return
	}

	if r.Method == http.MethodDelete {
		id := p.getId(rw, r)
		p.deleteProduct(rw, r, id)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProduct()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := p.handleData(rw, r)

	data.AddProduct(prod)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request, id int) {
	p.l.Println("Handle PUT Product")

	prod := p.handleData(rw, r)

	err := data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

func (p *Products) deleteProduct(rw http.ResponseWriter, r *http.Request, id int) {
	p.l.Println("Handle DELETE Product")

	err := data.DeleteProduct(id)
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

}

func (p *Products) handleData(rw http.ResponseWriter, r *http.Request) *data.Product {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	return prod
}

func (p *Products) getId(rw http.ResponseWriter, r *http.Request) int {
	reg := regexp.MustCompile(`/([0-9]+)`)
	g := reg.FindAllStringSubmatch(r.URL.Path, -1)

	if len(g) != 1 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
	}

	if len(g[0]) != 2 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
	}

	if id <= 0 {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	p.l.Println("got id:", id)

	return id
}
