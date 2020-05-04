package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// interface
// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	// handle the request for a list of products
// 	if r.Method == http.MethodGet {
// 		p.getProducts(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		p.addProduct(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPut {
// 		p.l.Println("PUT", r.URL.Path)
// 		// expect the id in the URI
// 		reg := regexp.MustCompile(`/([0-9]+)`)
// 		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

// 		if len(g) != 1 {
// 			p.l.Println("Invalid URI more than one id")
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		if len(g[0]) != 2 {
// 			p.l.Println("Invalid URI more than one capture group")
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		idString := g[0][1]
// 		id, err := strconv.Atoi(idString)
// 		if err != nil {
// 			p.l.Println("Invalid URI unable to convert to numer", idString)
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		p.updateProducts(id, rw, r)
// 		return
// 	}

// 	// catch all
// 	// if no method is satisfied return an error
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

// getProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	// prod := &data.Product{}

	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	// }

	// unmarshalling and creating the Product has been moved into middleware
	// here we are retrieving the product from the context and casting it
	// to a `data.Product`
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // `id` was previously passed into UpdateProducts
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product")
	// prod := &data.Product{}

	// err = prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	// }

	// unmarshalling and creating the Product has been moved into middleware
	// here we are retrieving the product from the context and casting it
	// to a `data.Product`
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// strings can be used as context keys, however using types is preferable
type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// Add product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		// call the next handler, which can be another middleware in the chain or the final handler
		next.ServeHTTP(rw, req)

	})
}
