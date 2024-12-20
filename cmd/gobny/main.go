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
	router.HandleFunc("POST /users", handler.CreateUser)
	router.HandleFunc("GET /users/{id}/carts", handler.GetUserCart)
	router.HandleFunc("POST /carts", handler.CreateCart)
	router.HandleFunc("POST /carts/items", handler.CreateCartItem)
	router.HandleFunc("GET /carts/{id}/items", handler.GetCartItems)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Server started on port 3000!")
	log.Fatal(server.ListenAndServe())
}
