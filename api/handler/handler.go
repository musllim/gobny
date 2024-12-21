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
	database "github.com/musllim/gobny/internal/database/gobny"
	"golang.org/x/crypto/bcrypt"
)

func throwError(w http.ResponseWriter, message string, status int, err error) {
	w.WriteHeader(status)
	fmt.Println(message, err.Error())
	w.Write([]byte(message))
}
func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	products, err := queries.GetProducts(ctx)

	if err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)

	cartId := r.PathValue("id")
	id, err := strconv.Atoi(cartId)

	if err != nil {
		throwError(w, "Provide a valid cart id", http.StatusBadRequest, err)
		return
	}

	products, err := queries.DeleteProduct(ctx, int32(id))

	if err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}
func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)

	cartId := r.PathValue("id")
	id, err := strconv.Atoi(cartId)

	if err != nil {
		throwError(w, "Provide a valid cart id", http.StatusBadRequest, err)
		return
	}

	products, err := queries.GetProduct(ctx, int32(id))

	if err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
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
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	var product database.CreateProductParams

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
		return
	}

	if _, err := queries.CreateProduct(ctx, product); err != nil {
		throwError(w, "products query failed", http.StatusInternalServerError, err)
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
		throwError(w, "db connection failed", http.StatusInternalServerError, err)
		return
	}

	type UserParams struct {
		Email    string
		Password string
	}
	defer conn.Close(ctx)
	var user UserParams

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		throwError(w, "Users Query failed", http.StatusInternalServerError, err)
		return
	}

	users, err := queries.GetUser(ctx, user.Email)

	if err != nil {
		throwError(w, "Users Query failed", http.StatusInternalServerError, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(user.Password)); err != nil {
		throwError(w, "Invalid creadentials", http.StatusUnauthorized, err)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(users.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		throwError(w, "Faied to generate token", http.StatusInternalServerError, err)
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
		throwError(w, "User Query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	var user database.CreateUserParams

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		throwError(w, "Make sure to send corect request body", http.StatusInternalServerError, err)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	user.Password = string(hashed)

	if _, err := queries.CreateUser(ctx, user); err != nil {
		throwError(w, "User creation failed:", http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))

}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "User Query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	id, err := strconv.Atoi(r.URL.RawQuery)

	if err != nil {
		throwError(w, "Provide a valid user id:", http.StatusBadRequest, err)
		return
	}

	user, err := queries.GetUserById(ctx, int32(id))

	if err != nil {
		throwError(w, "User fetch failed", http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user.Password = ""
	if err := json.NewEncoder(w).Encode(user); err != nil {
		throwError(w, "failed to get cart failed", http.StatusInternalServerError, err)
		return
	}
}
func CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "Products Query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	id, err := strconv.Atoi(r.URL.RawQuery)

	if err != nil {
		throwError(w, "Provide a valid user id:", http.StatusBadRequest, err)
		return
	}

	if _, err := queries.CreateUserCart(ctx, int32(id)); err != nil {
		throwError(w, "Cart creation failed", http.StatusInternalServerError, err)
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
		throwError(w, "Cart item Query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	var cartItem database.CreateCartItemParams

	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		throwError(w, "Make sure to send corect request body", http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.Atoi(r.URL.RawQuery)

	if err != nil {
		throwError(w, "Provide a valid user id:", http.StatusBadRequest, err)
		return
	}

	cart, err := queries.GetUserCart(ctx, int32(id))

	cartItem.Cartid = int32(cart.ID)

	if _, err := queries.CreateCartItem(ctx, cartItem); err != nil {
		throwError(w, "Cart item creation failed", http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cart item created successfully"))

}

func GetUserCart(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "User Query failed", http.StatusBadRequest, err)
		return
	}

	defer conn.Close(ctx)
	id, err := strconv.Atoi(r.URL.RawQuery)

	if err != nil {
		throwError(w, "Provide a valid user id", http.StatusBadRequest, err)
		return
	}

	cart, err := queries.GetUserCart(ctx, int32(id))

	if err != nil {
		throwError(w, "failed to get cart failed", http.StatusBadRequest, err)
		return
	}

	if err := json.NewEncoder(w).Encode(cart); err != nil {
		throwError(w, "failed to get cart failed", http.StatusInternalServerError, err)
		return
	}
}

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	queries := database.New(conn)

	if err != nil {
		throwError(w, "User Query failed", http.StatusInternalServerError, err)
		return
	}

	defer conn.Close(ctx)
	cartId := r.PathValue("id")
	id, err := strconv.Atoi(cartId)

	if err != nil {
		throwError(w, "Provide a valid cart id", http.StatusBadRequest, err)
		return
	}

	cart, err := queries.GetCartItems(ctx, int32(id))

	if err != nil {
		throwError(w, "failed to get cart failed", http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(cart); err != nil {
		throwError(w, "failed to get cart failed", http.StatusInternalServerError, err)
		return
	}
}
