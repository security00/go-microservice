package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/security00/go-microservice/product-api/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProduct(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	delRouter := sm.Methods(http.MethodDelete).Subrouter()
	delRouter.HandleFunc("/{id[0-9]+}", ph.DeleteProduct)

	s := &http.Server{
		Addr:         ":9000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			l.Fatal(err)
		}
	}()

	l.Println("Starting Server, Listening Port 9000...")

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	defer close(sigChan)
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	ct, _ := context.WithTimeout(context.Background(), 30*time.Second)

	_ = s.Shutdown(ct)
}
