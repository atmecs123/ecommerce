package middleware

import (
	"encoding/json"
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

		next.ServeHTTP(w, r)
	})
}

func generateJWTToken() (string, error) {
	//create the claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
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
