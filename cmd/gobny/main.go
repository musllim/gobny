package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/musllim/gobny/api/handler"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /products", handler.GetProducts)
	router.HandleFunc("POST /products", handler.CreateProducts)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Server started on port 3000!")
	log.Fatal(server.ListenAndServe())
}
