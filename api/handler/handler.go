package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Println("db connection failed:", err.Error())
		w.Write([]byte("User Query failed"))

		return
	}
	type UserParams struct {
		Email    string
		Password string
	}
	defer conn.Close(ctx)
	var user UserParams

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Users Query failed:", err.Error())
		w.Write([]byte("Make sure to send corect request body"))
		return
	}

	users, err := queries.GetUser(ctx, pgtype.Text{String: user.Email, Valid: true})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("User query failed:", err.Error())
		w.Write([]byte("User query failed"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password.String), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid creadentials"))
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(users.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, error := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if error != nil {
		fmt.Println(error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Faied to generate token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))

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

func AddItemToCart(w http.ResponseWriter, r *http.Request) {
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

func GetUserCart(w http.ResponseWriter, r *http.Request) {
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
	userId := r.PathValue("id")
	id, error := strconv.Atoi(userId)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Provide a valid user id:", error.Error())
		w.Write([]byte("Provide a valid user id"))
		return
	}

	cart, error := queries.GetUserCart(ctx, pgtype.Int4{Int32: int32(id), Valid: true})

	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("failed to get cart failed:", error.Error())
		w.Write([]byte("failed to get cart failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(cart); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("failed to get cart failed:", err.Error())
		w.Write([]byte("failed to get cart failed"))
		return
	}
}

func GetCartItems(w http.ResponseWriter, r *http.Request) {
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
	cartId := r.PathValue("id")
	id, error := strconv.Atoi(cartId)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Provide a valid user id:", error.Error())
		w.Write([]byte("Provide a valid user id"))
		return
	}

	cart, error := queries.GetCartItems(ctx, pgtype.Int4{Int32: int32(id), Valid: true})

	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("failed to get cart failed:", error.Error())
		w.Write([]byte("failed to get cart failed"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(cart); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("failed to get cart failed:", err.Error())
		w.Write([]byte("failed to get cart failed"))
		return
	}
}
