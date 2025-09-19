package api

import (
	"crudl_service/src/config"
	"crudl_service/src/service"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var request types.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := service.AuthenticateUser(request.Username, request.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(userID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	sendAuthResponse(w, token, userID)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request types.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	userID, err := service.CreateUser(request.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	token, err := generateJWT(userID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	sendAuthResponse(w, token, userID)
}

func sendAuthResponse(w http.ResponseWriter, token, userID string) {
	response := types.AuthResponse{
		Token:  token,
		UserID: userID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateJWT(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.CurrentConfig.JWT.SecretKey))
}

func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("User-ID", claims.UserID)
		next(w, r)
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}

	return authHeader
}

func parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.CurrentConfig.JWT.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func CheckTaskOwner(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("User-ID")
		taskID := r.URL.Query().Get("task_id")

		if taskID == "" {
			http.Error(w, "Missing task_id", http.StatusBadRequest)
			return
		}

		isOwner, err := service.CheckTaskOwnership(userID, taskID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !isOwner {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
