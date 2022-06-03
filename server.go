package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yuhengfdada/go-microservices/handlers"
)

func main() {
	l := log.New(os.Stdout, "server-api ", log.LstdFlags)
	hello := handlers.NewHello(l)
	product := handlers.NewProduct(l)

	sm := http.NewServeMux()
	sm.Handle("/", hello)
	sm.Handle("/products", product)

	go func() {
		err := http.ListenAndServe(":8080", sm)
		if err != nil {
			log.Fatalf("error!")
			return
		}
	}()
	l.Println("starting server on port 8080")
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
}
