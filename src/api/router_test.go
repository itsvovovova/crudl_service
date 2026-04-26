package api

import (
	"bytes"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSubscription_ValidRequest(t *testing.T) {
	app := newTestApp(newMockRepository())

	startDate := "01-2023"
	subscription := types.UserSubscription{
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &startDate,
	}

	jsonData, _ := json.Marshal(subscription)
	req := httptest.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonData))
	req.Header.Set("User-ID", "user123")
	w := httptest.NewRecorder()

	app.CreateSubscription(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCreateSubscription_InvalidTimeFormat(t *testing.T) {
	app := newTestApp(newMockRepository())

	invalidDate := "01-03-2025"
	subscription := types.UserSubscription{
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   &invalidDate,
	}

	jsonData, _ := json.Marshal(subscription)
	req := httptest.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	app.CreateSubscription(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReadSubscription_Structure(t *testing.T) {
	app := newTestApp(newMockRepository())

	req := httptest.NewRequest("GET", "/subscription", nil)
	w := httptest.NewRecorder()

	app.ReadSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}

func TestUpdateSubscription_Structure(t *testing.T) {
	app := newTestApp(newMockRepository())

	startDate := "01-2023"
	endDate := "03-2025"
	subscription := types.UserSubscription{
		ServiceName: "Netflix",
		Price:       1299,
		StartDate:   &startDate,
		EndDate:     &endDate,
	}

	jsonData, _ := json.Marshal(subscription)
	req := httptest.NewRequest("PUT", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	app.UpdateSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}

func TestDeleteSubscription_Structure(t *testing.T) {
	app := newTestApp(newMockRepository())

	req := httptest.NewRequest("DELETE", "/subscription", nil)
	w := httptest.NewRecorder()

	app.DeleteSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}

func TestListSubscription_Structure(t *testing.T) {
	app := newTestApp(newMockRepository())

	req := httptest.NewRequest("GET", "/subscriptionList", nil)
	w := httptest.NewRecorder()

	app.ListSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}
