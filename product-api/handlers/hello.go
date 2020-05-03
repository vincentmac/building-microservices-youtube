package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Running Hello Handler")

	// read the body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.l.Println("Error reading body", err)

		http.Error(rw, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// write the response
	fmt.Fprintf(rw, "Hello %s", b)
}
