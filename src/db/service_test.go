package db

import (
	"crudl_service/src/types"
	"testing"
)

func newNilRepo() *postgresRepository {
	return &postgresRepository{db: nil}
}

func TestCreate_NilDB(t *testing.T) {
	r := newNilRepo()
	startDate := "01-2023"
	_, err := r.Create(&types.UserSubscription{
		ServiceName: "Netflix", Price: 999, UserId: "user123", StartDate: &startDate,
	})
	if err == nil {
		t.Error("Expected error with nil db")
	}
}

func TestGet_NilDB(t *testing.T) {
	r := newNilRepo()
	result, err := r.Get(1)
	if err == nil {
		t.Error("Expected error with nil db")
	}
	if result != nil {
		t.Error("Expected nil result")
	}
}

func TestUpdate_NilDB(t *testing.T) {
	r := newNilRepo()
	startDate := "01-2023"
	err := r.Update(&types.UserSubscription{
		ServiceName: "Netflix", Price: 999, UserId: "user123", StartDate: &startDate,
	})
	if err == nil {
		t.Error("Expected error with nil db")
	}
}

func TestDelete_NilDB(t *testing.T) {
	r := newNilRepo()
	if err := r.Delete(1); err == nil {
		t.Error("Expected error with nil db")
	}
}

func TestList_NilDB(t *testing.T) {
	r := newNilRepo()
	result, err := r.List("user123", nil, 10)
	if err == nil {
		t.Error("Expected error with nil db")
	}
	if result != nil {
		t.Error("Expected nil result")
	}
}

func TestSum_NilDB(t *testing.T) {
	r := newNilRepo()
	result, err := r.Sum(&types.UserSumSubscriptionRequest{
		UserId: "user123", StartDate: "01-2023", EndDate: "12-2023",
	})
	if err == nil {
		t.Error("Expected error with nil db")
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestSum_NilData(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Function should not panic: %v", r)
		}
	}()
	r := newNilRepo()
	_, err := r.Sum(nil)
	if err == nil {
		t.Error("Expected error for nil data")
	}
}
