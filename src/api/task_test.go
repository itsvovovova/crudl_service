package api

import (
	"context"
	"crudl_service/src/db"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockRepository struct {
	subscriptions map[int64]*types.UserSubscription
	nextID        int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		subscriptions: make(map[int64]*types.UserSubscription),
		nextID:        1,
	}
}

func (m *mockRepository) Create(data *types.UserSubscription) (int64, error) {
	data.Id = m.nextID
	m.subscriptions[m.nextID] = data
	m.nextID++
	return data.Id, nil
}

func (m *mockRepository) Get(id int64) (*types.UserSubscription, error) {
	if sub, exists := m.subscriptions[id]; exists {
		return sub, nil
	}
	return nil, &db.NotFoundError{}
}

func (m *mockRepository) Update(data *types.UserSubscription) error {
	if _, exists := m.subscriptions[data.Id]; !exists {
		return &db.NotFoundError{}
	}
	m.subscriptions[data.Id] = data
	return nil
}

func (m *mockRepository) Delete(id int64) error {
	if _, exists := m.subscriptions[id]; !exists {
		return &db.NotFoundError{}
	}
	delete(m.subscriptions, id)
	return nil
}

func (m *mockRepository) List(userID string, afterID *int64, limit int) ([]types.UserSubscription, error) {
	var result []types.UserSubscription
	count := 0

	for _, sub := range m.subscriptions {
		if sub.UserId == userID {
			if afterID == nil || sub.Id > *afterID {
				result = append(result, *sub)
				count++
				if limit > 0 && count >= limit {
					break
				}
			}
		}
	}

	return result, nil
}

func (m *mockRepository) Sum(data *types.UserSumSubscriptionRequest) (int64, error) {
	var sum int64
	for _, sub := range m.subscriptions {
		if sub.UserId == data.UserId {
			sum += sub.Price
		}
	}
	return sum, nil
}

func TestReadSubscription_ValidID(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	startDate := "01-2023"
	subscription := &types.UserSubscription{
		Id:          1,
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &startDate,
	}

	repo.subscriptions[1] = subscription

	req := httptest.NewRequest("GET", "/subscription/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	ReadSubscription(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response types.UserSubscription
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal response")
	}

	if response.ServiceName != "Netflix" {
		t.Errorf("Expected service name 'Netflix', got '%s'", response.ServiceName)
	}
}

func TestReadSubscription_InvalidID(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	req := httptest.NewRequest("GET", "/subscription/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	ReadSubscription(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReadSubscription_NotFound(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	req := httptest.NewRequest("GET", "/subscription/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	ReadSubscription(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListSubscription_ValidUser(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	startDate := "01-2023"
	subscription1 := &types.UserSubscription{
		Id:          1,
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &startDate,
	}
	subscription2 := &types.UserSubscription{
		Id:          2,
		ServiceName: "Spotify",
		Price:       499,
		UserId:      "user123",
		StartDate:   &startDate,
	}

	repo.subscriptions[1] = subscription1
	repo.subscriptions[2] = subscription2

	req := httptest.NewRequest("GET", "/subscriptionList?user_id=user123", nil)
	w := httptest.NewRecorder()

	ListSubscription(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Data        []types.UserSubscription `json:"data"`
		NextAfterID *int64                   `json:"next_after_id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal response")
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", len(response.Data))
	}
}

func TestListSubscription_WithLimit(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	startDate := "01-2023"
	for i := 1; i <= 5; i++ {
		subscription := &types.UserSubscription{
			Id:          int64(i),
			ServiceName: "Service" + string(rune(i)),
			Price:       999,
			UserId:      "user123",
			StartDate:   &startDate,
		}
		repo.subscriptions[int64(i)] = subscription
	}

	req := httptest.NewRequest("GET", "/subscriptionList?user_id=user123&limit=3", nil)
	w := httptest.NewRecorder()

	ListSubscription(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Data        []types.UserSubscription `json:"data"`
		NextAfterID *int64                   `json:"next_after_id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal response")
	}

	if len(response.Data) != 3 {
		t.Errorf("Expected 3 subscriptions, got %d", len(response.Data))
	}
}

func TestDeleteSubscription_ValidID(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	startDate := "01-2023"
	subscription := &types.UserSubscription{
		Id:          1,
		ServiceName: "Netflix",
		Price:       999,
		UserId:      "user123",
		StartDate:   &startDate,
	}

	repo.subscriptions[1] = subscription

	req := httptest.NewRequest("DELETE", "/subscription/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	DeleteSubscription(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if _, exists := repo.subscriptions[1]; exists {
		t.Error("Subscription should be deleted")
	}
}

func TestDeleteSubscription_NotFound(t *testing.T) {
	repo := newMockRepository()
	CurrentRepository = repo

	req := httptest.NewRequest("DELETE", "/subscription/999", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	DeleteSubscription(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
