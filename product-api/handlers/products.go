package handlers

import (
	"log"
	"net/http"

	"github.com/vincentmac/building-microservices-youtube/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// fulfill handler interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle get - curl localhost:9090 | jq
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// handle update
	// if r.Method == http.MethodPut {

	// }

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed) // curl localhost:9090 -XDELETE -v| jq
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// d, err := json.Marshal(lp) // method 1 to return JSON. This allocates space for `d`
	err := lp.ToJSON(rw) // method 2 for JSON. This is quicker and preferred.
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	// rw.Write(d) // part of method 1
}
