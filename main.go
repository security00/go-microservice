package main

import (
	"context"
	"github.com/security00/go-microservice/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProduct(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

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
