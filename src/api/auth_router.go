package api

import (
	"crudl_service/src/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (a *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request types.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := a.repo.GetUserByUsername(request.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := a.generateJWT(user.ID)
	if err != nil {
		log.WithError(err).Error("Failed to generate JWT")
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	a.sendAuthResponse(w, token, user.ID, http.StatusOK)
}

func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request types.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("Failed to hash password")
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	userID, err := a.repo.CreateUser(request.Username, string(hashedPassword))
	if err != nil {
		log.WithError(err).Error("Failed to create user")
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	token, err := a.generateJWT(userID)
	if err != nil {
		log.WithError(err).Error("Failed to generate JWT")
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	a.sendAuthResponse(w, token, userID, http.StatusCreated)
}

func (a *App) sendAuthResponse(w http.ResponseWriter, token, userID string, status int) {
	body, err := json.Marshal(types.AuthResponse{Token: token, UserID: userID})
	if err != nil {
		log.WithError(err).Error("Failed to marshal auth response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (a *App) generateJWT(userID string) (string, error) {
	if a.jwtSecret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func (a *App) ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		claims, err := a.parseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		r.Header.Set("User-ID", claims.UserID)
		next(w, r)
	}
}

func (a *App) parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(a.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}
	return claims, nil
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
