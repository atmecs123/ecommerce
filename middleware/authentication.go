package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Credentials struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	username = "admin"
	password = "admin"
)

// Secret key used to sign tokens
var jwtKey = []byte("my_secret_key")

// /The middleware will check the Authorization header for a valid JWT.
// If the token is missing or invalid, it will block the request; otherwise, it will pass control to the next handler.
func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, "uuid", "123456789")
		//The middleware sets a timeout for the request context.
		//If the handler takes longer than the specified timeout duration, the request will be aborted.
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		fmt.Println("#### inside authmiddleware ####")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be bearer token", http.StatusUnauthorized)
			return
		}
		tokenString := authHeaderParts[1]
		claims := &jwt.StandardClaims{}
		// Parse the JWT string and store the result in `claims`.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil && !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// Create a channel to signal the completion of the request handling
		done := make(chan bool, 1)
		//It runs the next handler in a separate goroutine,
		//which allows the middleware to listen for the completion of the handler or the timeout signal, whichever comes first.
		go func() {
			next.ServeHTTP(w, r)
			done <- true
		}()

		// Listen on multiple channels using select
		select {
		case <-ctx.Done(): // If the context's deadline is exceeded
			switch ctx.Err() {
			case context.DeadlineExceeded:
				w.WriteHeader(http.StatusGatewayTimeout)
				fmt.Fprintln(w, "Request timed out")
			case context.Canceled:
				fmt.Fprintln(w, "Request was canceled")
			}
		//A channel (done) is used to signal the completion of the handler.
		//This allows the middleware to stop waiting if the handler finishes before the timeout.
		case <-done:

		}
	})
}

func generateJWTToken() (string, error) {
	//create the claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(20 * time.Minute).Unix(),
		Issuer:    "ecommerce",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign the token with secret key
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func Authentication(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Not a valid request", http.StatusBadRequest)
		return
	}
	// TODO: Validate credentials here. For now, we're using dummy data.
	if cred.Username != username || cred.Password != password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := generateJWTToken()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
