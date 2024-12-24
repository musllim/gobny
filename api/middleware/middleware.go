package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}
func IsAutenticated(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		fmt.Println("Authorization", authorization)
		if !strings.HasPrefix(authorization, "Bearer ") {
			writeUnauthed(w)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		jwtToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {

			writeUnauthed(w)
			return
		}

		claims, ok := jwtToken.Claims.(*jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("failed to parse claims"))
			return
		}
		expDate, err := claims.GetExpirationTime()
		if err != nil {
			writeUnauthed(w)
			return
		}
		if time.Now().Unix() > expDate.Unix() {
			writeUnauthed(w)
			return
		}
		subject, error := claims.GetSubject()
		if error != nil {
			writeUnauthed(w)
			return
		}
		r.URL.RawQuery = subject
		next.ServeHTTP(w, r)
	})
}

func IsOwnCart(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
