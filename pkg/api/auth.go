package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type SigninRequest struct {
	Password string `json:"password"`
}

var jwtKey = []byte("todo-secret-key")

func passwordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func createToken(password string) (string, error) {
	claims := jwt.MapClaims{
		"hash": passwordHash(password),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateToken(tokenString string, password string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hash, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	return hash == passwordHash(password)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	if req.Password != pass {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "неверный пароль",
		})
		return
	}

	token, err := createToken(pass)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusBadRequest, map[string]string{
		"token": token,
	})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if pass != "" {
			cookie, err := r.Cookie("token")

			if err != nil || !validateToken(cookie.Value, pass) {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	}
}
