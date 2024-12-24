package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	database "github.com/musllim/gobny/internal/database/gobny"
	"github.com/musllim/gobny/pkg"
	"golang.org/x/crypto/bcrypt"
)

func GetProducts(ctx context.Context, queries *database.Queries) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products, err := queries.GetProducts(ctx)

		if err != nil {
			pkg.ThrowError(w, "products query failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	})

}

func DeleteProduct(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cartId := r.PathValue("id")
		id, err := strconv.Atoi(cartId)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid cart id", http.StatusBadRequest, err)
			return
		}

		products, err := queries.DeleteProduct(ctx, int32(id))

		if err != nil {
			pkg.ThrowError(w, "products query failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	})
}

func GetProduct(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cartId := r.PathValue("id")
		id, err := strconv.Atoi(cartId)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid cart id", http.StatusBadRequest, err)
			return
		}

		products, err := queries.GetProduct(ctx, int32(id))

		if err != nil {
			pkg.ThrowError(w, "products query failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	})
}

func CreateProducts(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var product database.CreateProductParams

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			pkg.ThrowError(w, "products query failed", http.StatusInternalServerError, err)
			return
		}

		if _, err := queries.CreateProduct(ctx, product); err != nil {
			pkg.ThrowError(w, "products query failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Products created successfully"))

	})
}

func LoginUser(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type UserParams struct {
			Email    string
			Password string
		}
		var user UserParams

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			pkg.ThrowError(w, "Users Query failed", http.StatusInternalServerError, err)
			return
		}

		users, err := queries.GetUser(ctx, user.Email)

		if err != nil {
			pkg.ThrowError(w, "Users Query failed", http.StatusInternalServerError, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(user.Password)); err != nil {
			pkg.ThrowError(w, "Invalid creadentials", http.StatusUnauthorized, err)
			return
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": strconv.Itoa(int(users.ID)),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			pkg.ThrowError(w, "Faied to generate token", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	})
}

func CreateUser(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var user database.CreateUserParams

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			pkg.ThrowError(w, "Make sure to send corect request body", http.StatusInternalServerError, err)
			return
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

		user.Password = string(hashed)

		if _, err := queries.CreateUser(ctx, user); err != nil {
			pkg.ThrowError(w, "User creation failed:", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User created successfully"))

	})
}

func UserProfile(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.RawQuery)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid user id:", http.StatusBadRequest, err)
			return
		}

		user, err := queries.GetUserById(ctx, int32(id))

		if err != nil {
			pkg.ThrowError(w, "User fetch failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		user.Password = ""
		if err := json.NewEncoder(w).Encode(user); err != nil {
			pkg.ThrowError(w, "failed to get cart failed", http.StatusInternalServerError, err)
			return
		}
	})
}

func CreateCart(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.URL.RawQuery)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid user id:", http.StatusBadRequest, err)
			return
		}

		if _, err := queries.CreateUserCart(ctx, int32(id)); err != nil {
			pkg.ThrowError(w, "Cart creation failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cart created successfully"))

	})
}

func CreateCartItem(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cartItem database.CreateCartItemParams

		if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
			pkg.ThrowError(w, "Make sure to send corect request body", http.StatusInternalServerError, err)
			return
		}
		id, err := strconv.Atoi(r.URL.RawQuery)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid user id:", http.StatusBadRequest, err)
			return
		}

		cart, err := queries.GetUserCart(ctx, int32(id))

		cartItem.Cartid = int32(cart.ID)

		if _, err := queries.CreateCartItem(ctx, cartItem); err != nil {
			pkg.ThrowError(w, "Cart item creation failed", http.StatusBadRequest, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cart item created successfully"))

	})
}
func GetUserCart(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.RawQuery)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid user id", http.StatusBadRequest, err)
			return
		}

		cart, err := queries.GetUserCart(ctx, int32(id))

		if err != nil {
			pkg.ThrowError(w, "failed to get cart failed", http.StatusBadRequest, err)
			return
		}

		if err := json.NewEncoder(w).Encode(cart); err != nil {
			pkg.ThrowError(w, "failed to get cart failed", http.StatusInternalServerError, err)
			return
		}
	})
}
func GetCartItems(ctx context.Context, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cartId := r.PathValue("id")
		id, err := strconv.Atoi(cartId)

		if err != nil {
			pkg.ThrowError(w, "Provide a valid cart id", http.StatusBadRequest, err)
			return
		}

		cart, err := queries.GetCartItems(ctx, int32(id))

		if err != nil {
			pkg.ThrowError(w, "failed to get cart failed", http.StatusInternalServerError, err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(cart); err != nil {
			pkg.ThrowError(w, "failed to get cart failed", http.StatusInternalServerError, err)
			return
		}
	})
}
