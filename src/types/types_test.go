package types

import (
	"encoding/json"
	"testing"
)

func TestUserSubscription_JSONMarshaling(t *testing.T) {
	startDate := "01-2023"
	endDate := "12-2023"

	subscription := UserSubscription{
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &startDate,
		EndDate:     &endDate,
	}

	jsonData, err := json.Marshal(subscription)
	if err != nil {
		t.Fatalf("Failed to marshal UserSubscription: %v", err)
	}

	var unmarshaled UserSubscription
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserSubscription: %v", err)
	}

	if unmarshaled.ServiceName != subscription.ServiceName {
		t.Errorf("Expected service name '%s', got '%s'", subscription.ServiceName, unmarshaled.ServiceName)
	}
	if unmarshaled.Price != subscription.Price {
		t.Errorf("Expected price %d, got %d", subscription.Price, unmarshaled.Price)
	}
	if unmarshaled.UserId != subscription.UserId {
		t.Errorf("Expected user id '%s', got '%s'", subscription.UserId, unmarshaled.UserId)
	}
}

func TestUserSubscriptionData_JSONMarshaling(t *testing.T) {
	data := UserSubscriptionData{
		UserId:      "user123",
		ServiceName: "Spotify",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal UserSubscriptionData: %v", err)
	}

	var unmarshaled UserSubscriptionData
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserSubscriptionData: %v", err)
	}

	if unmarshaled.UserId != data.UserId {
		t.Errorf("Expected user id '%s', got '%s'", data.UserId, unmarshaled.UserId)
	}
	if unmarshaled.ServiceName != data.ServiceName {
		t.Errorf("Expected service name '%s', got '%s'", data.ServiceName, unmarshaled.ServiceName)
	}
}

func TestUserSumSubscriptionRequest_JSONMarshaling(t *testing.T) {
	request := UserSumSubscriptionRequest{
		UserId:    "user123",
		StartDate: "01-2023",
		EndDate:   "12-2023",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal UserSumSubscriptionRequest: %v", err)
	}

	var unmarshaled UserSumSubscriptionRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserSumSubscriptionRequest: %v", err)
	}

	if unmarshaled.UserId != request.UserId {
		t.Errorf("Expected user id '%s', got '%s'", request.UserId, unmarshaled.UserId)
	}
}

func TestUserSubscriptionSumResponse_JSONMarshaling(t *testing.T) {
	response := UserSubscriptionSumResponse{
		UserId:     "user123",
		CurrentSum: 2500,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal UserSubscriptionSumResponse: %v", err)
	}

	var unmarshaled UserSubscriptionSumResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserSubscriptionSumResponse: %v", err)
	}

	if unmarshaled.UserId != response.UserId {
		t.Errorf("Expected user id '%s', got '%s'", response.UserId, unmarshaled.UserId)
	}
	if unmarshaled.CurrentSum != response.CurrentSum {
		t.Errorf("Expected current sum %d, got %d", response.CurrentSum, unmarshaled.CurrentSum)
	}
}
