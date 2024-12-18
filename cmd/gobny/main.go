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
	router.HandleFunc("POST /cart", handler.CreateCart)
	router.HandleFunc("POST /cart/items", handler.CreateCartItem)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Server started on port 3000!")
	log.Fatal(server.ListenAndServe())
}
