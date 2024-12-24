package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/musllim/gobny/api/handler"
	"github.com/musllim/gobny/api/middleware"
	database "github.com/musllim/gobny/internal/database/gobny"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		log.Fatal("Db connection failed")
		return
	}

	defer conn.Close(ctx)

	authorizedRouter := http.NewServeMux()
	router := http.NewServeMux()

	router.Handle("/", middleware.IsAutenticated(authorizedRouter))

	router.HandleFunc("GET /products", handler.GetProducts(ctx, queries).ServeHTTP)
	router.HandleFunc("GET /products/{id}", handler.GetProduct(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("POST /products", handler.CreateProducts(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("GET /profile", handler.UserProfile(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("DELETE /products/{id}", handler.DeleteProduct(ctx, queries).ServeHTTP)
	router.HandleFunc("POST /users", handler.CreateUser(ctx, queries).ServeHTTP)
	router.HandleFunc("POST /login", handler.LoginUser(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("GET /carts", handler.GetUserCart(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("POST /carts", handler.CreateCart(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("POST /carts/items", handler.CreateCartItem(ctx, queries).ServeHTTP)
	authorizedRouter.HandleFunc("GET /carts/{id}/items", handler.GetCartItems(ctx, queries).ServeHTTP)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Server started on port 3000!")
	log.Fatal(server.ListenAndServe())
}
