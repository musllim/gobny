package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	database "github.com/musllim/gobny/internal/database/gobny"
	"golang.org/x/crypto/bcrypt"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println("db connection failed:", err.Error())
		w.Write([]byte("Products Query failed"))

		return
	}

	defer conn.Close(ctx)
	products, err := queries.GetProducts(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Products Query failed:", err.Error())
		w.Write([]byte("Products Query failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

func CreateProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println("db connection failed:", err.Error())
		w.Write([]byte("Products Query failed"))

		return
	}

	defer conn.Close(ctx)
	var product database.CreateProductParams

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Products Query failed:", err.Error())
		w.Write([]byte("Make sure to send corect request body"))
		return
	}

	products, err := queries.CreateProduct(ctx, product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Products creation failed:", err.Error(), products)
		w.Write([]byte("Products creation failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Products created successfully"))

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println("db connection failed:", err.Error())
		w.Write([]byte("User Query failed"))

		return
	}

	defer conn.Close(ctx)
	var user database.CreateUserParams

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Users Query failed:", err.Error())
		w.Write([]byte("Make sure to send corect request body"))
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password.String), 8)

	user.Password.String = string(hashed)
	users, err := queries.CreateUser(ctx, user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("User creation failed:", err.Error(), users)
		w.Write([]byte("User creation failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))

}

func CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println("db connection failed:", err.Error())
		w.Write([]byte("Products Query failed"))

		return
	}

	defer conn.Close(ctx)
	var cart database.Cart

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Cart creation failed:", err.Error())
		w.Write([]byte("Make sure to send corect request body"))
		return
	}

	carts, err := queries.CreateUserCart(ctx, cart.UserID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Cart creation failed:", err.Error(), carts)
		w.Write([]byte("Cart creation failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cart created successfully"))

}

func CreateCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println("db connection failed:", err.Error())
		w.Write([]byte("Cart item Query failed"))

		return
	}

	defer conn.Close(ctx)
	var cartItem database.CreateCartItemParams

	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Cart item creation failed:", err.Error())
		w.Write([]byte("Make sure to send corect request body"))
		return
	}

	carts, err := queries.CreateCartItem(ctx, cartItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Cart item creation failed:", err.Error(), carts)
		w.Write([]byte("Cart item creation failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cart item created successfully"))

}
