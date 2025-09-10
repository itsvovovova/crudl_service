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
	startDate := "01-2023"
	subscription := types.UserSubscription{
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &startDate,
		EndDate:     nil,
	}

	jsonData, _ := json.Marshal(subscription)
	req := httptest.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	CreateSubscription(w, req)

	if w.Code == http.StatusBadRequest &&
		w.Body.String() == "Incorrect time format" {
		t.Error("Time format validation failed for valid date")
	}
}

func TestCreateSubscription_InvalidTimeFormat(t *testing.T) {
	invalidDate := "01-03-2025"
	subscription := types.UserSubscription{
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &invalidDate,
		EndDate:     nil,
	}

	jsonData, _ := json.Marshal(subscription)
	req := httptest.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	CreateSubscription(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid time format, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReadSubscription_Structure(t *testing.T) {
	subscriptionData := types.UserSubscriptionData{
		UserId:      "user123",
		ServiceName: "Netflix",
	}

	jsonData, _ := json.Marshal(subscriptionData)
	req := httptest.NewRequest("GET", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	ReadSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}

func TestUpdateSubscription_Structure(t *testing.T) {
	startDate := "01-2023"
	subscription := types.UserSubscription{
		ServiceName: "Netflix",
		Price:       1299,
		UserId:      "user123",
		StartDate:   &startDate,
		EndDate:     nil,
	}

	jsonData, _ := json.Marshal(subscription)
	req := httptest.NewRequest("PUT", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	UpdateSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}

func TestDeleteSubscription_Structure(t *testing.T) {
	subscriptionData := types.UserSubscriptionData{
		UserId:      "user123",
		ServiceName: "Netflix",
	}

	jsonData, _ := json.Marshal(subscriptionData)
	req := httptest.NewRequest("DELETE", "/subscription", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	DeleteSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}

func TestListSubscription_Structure(t *testing.T) {
	userRequest := types.UserRequest{
		UserId: "user123",
	}

	jsonData, _ := json.Marshal(userRequest)
	req := httptest.NewRequest("GET", "/subscriptionList", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	ListSubscription(w, req)

	if w.Code == 0 {
		t.Error("Handler did not set any response code")
	}
}
