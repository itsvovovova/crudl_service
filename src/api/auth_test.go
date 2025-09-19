package api

import (
	"bytes"
	"crudl_service/src/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func setupTestConfig() {
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	config.InitConfig()
}

func TestLoginUser_InvalidRequest(t *testing.T) {
	setupTestConfig()

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	LoginUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterUser_InvalidRequest(t *testing.T) {
	setupTestConfig()

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	RegisterUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	setupTestConfig()

	token, _ := generateJWT("user123")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("User-ID")
		if userID != "user123" {
			t.Errorf("Expected userID 'user123', got '%s'", userID)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	setupTestConfig()

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler := ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	})

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestValidateJWT_MissingToken(t *testing.T) {
	setupTestConfig()

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	handler := ValidateJWT(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without token")
	})

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCheckTaskOwner_MissingTaskID(t *testing.T) {
	req := httptest.NewRequest("GET", "/task", nil)
	req.Header.Set("User-ID", "user123")
	w := httptest.NewRecorder()

	handler := CheckTaskOwner(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without task_id")
	})

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGenerateJWT(t *testing.T) {
	setupTestConfig()

	token, err := generateJWT("user123")
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.CurrentConfig.JWT.SecretKey), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Expected valid token")
	}

	if claims.UserID != "user123" {
		t.Errorf("Expected userID 'user123', got '%s'", claims.UserID)
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"Bearer token", "Bearer abc123", "abc123"},
		{"Token without Bearer", "abc123", "abc123"},
		{"Empty header", "", ""},
		{"Bearer with space", "Bearer  abc123", " abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			result := extractToken(req)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
