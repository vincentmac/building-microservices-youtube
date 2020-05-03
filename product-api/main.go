package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

	"github.com/nicholasjackson/building-microservices-youtube/product-api/handlers"
)

func main() {
	// reqeusts to the path /goodbye with be handled by this function
	// http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
	// 	log.Println("Goodbye World")
	// })

	// any other request will be handled by this function
	// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	// })

	// Listen for connections on all ip addresses (0.0.0.0)
	// port 9090
	// log.Println("Starting Server")
	// err := http.ListenAndServe(":9090", sm)
	// log.Fatal(err)

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// put listenAndServe in a go func so it doesn't block
	// the gracefull shutdown via signal
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// handle shutdown gracefully
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// block until we recieve a signal
	sig := <-sigChan
	l.Println("Received termintate, graceful shutdown", sig)
	// do additional cleanup here, if needed, then shutdown server

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
