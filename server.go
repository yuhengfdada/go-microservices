package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/yuhengfdada/go-microservices/handlers"
)

func main() {
	l := log.New(os.Stdout, "server-api ", log.LstdFlags)
	productHandler := handlers.NewProductHandler(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareProductConversion)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareProductConversion)

	// set up the server to look for the swagger file
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redocHandler := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", redocHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

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
