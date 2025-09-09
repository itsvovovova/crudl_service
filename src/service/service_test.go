package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestReadUserData_ValidJSON(t *testing.T) {
	jsonData := `{"name":"John","age":30}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonData))
	w := httptest.NewRecorder()

	var result TestStruct
	ReadUserData(w, req, &result)

	if result.Name != "John" {
		t.Errorf("Expected name 'John', got '%s'", result.Name)
	}
	if result.Age != 30 {
		t.Errorf("Expected age 30, got %d", result.Age)
	}
}

func TestReadUserData_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":"John","age":}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
	w := httptest.NewRecorder()

	var result TestStruct
	ReadUserData(w, req, &result)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReadUserData_NotPointer(t *testing.T) {
	jsonData := `{"name":"John","age":30}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonData))
	w := httptest.NewRecorder()

	var result TestStruct
	ReadUserData(w, req, result) // Passing value instead of pointer

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
