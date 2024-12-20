package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/musllim/gobny/api/handler"
	"github.com/musllim/gobny/api/middleware"
)

func main() {
	authorizedRouter := http.NewServeMux()
	router := http.NewServeMux()

	router.Handle("/", middleware.IsAutenticated(authorizedRouter))

	router.HandleFunc("GET /products", handler.GetProducts)
	authorizedRouter.HandleFunc("POST /products", handler.CreateProducts)
	router.HandleFunc("POST /users", handler.CreateUser)
	router.HandleFunc("POST /login", handler.LoginUser)
	authorizedRouter.HandleFunc("GET /carts", handler.GetUserCart)
	authorizedRouter.HandleFunc("POST /carts", handler.CreateCart)
	authorizedRouter.HandleFunc("POST /carts/items", handler.CreateCartItem)
	authorizedRouter.HandleFunc("GET /carts/{id}/items", handler.GetCartItems)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Server started on port 3000!")
	log.Fatal(server.ListenAndServe())
}
