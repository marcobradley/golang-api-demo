package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSortArrayReturnsSortedList(t *testing.T) {
	input := []int{5, 3, 8, 1, 2, 7}
	got := sortArray(input)
	want := []int{1, 2, 3, 5, 7, 8}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected sorted array %v, got %v", want, got)
	}
}

func TestSortArrayDoesNotMutateInput(t *testing.T) {
	input := []int{4, 1, 4, 2}
	original := append([]int(nil), input...)

	_ = sortArray(input)

	if !reflect.DeepEqual(input, original) {
		t.Fatalf("sortArray mutated input: got %v, want %v", input, original)
	}
}

func TestSortArrayHandlesEmptyAndSingleItem(t *testing.T) {
	empty := []int{}
	if got := sortArray(empty); len(got) != 0 {
		t.Fatalf("expected empty sorted list, got %v", got)
	}

	single := []int{42}
	if got := sortArray(single); !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("expected single-item sorted list, got %v", got)
	}
}

func TestQuicksortEndpointSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	body := []byte(`{"array":[5,3,8,1,2,7]}`)
	req, err := http.NewRequest(http.MethodPost, "/quicksort", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, w.Code)
	}

	var got map[string][]int
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	want := []int{1, 2, 3, 5, 7, 8}
	if !reflect.DeepEqual(got["sorted"], want) {
		t.Errorf("expected sorted array %v but got %v", want, got["sorted"])
	}
}

func TestQuicksortEndpointInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	body := []byte(`{"array":`)
	req, err := http.NewRequest(http.MethodPost, "/quicksort", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d but got %d", http.StatusBadRequest, w.Code)
	}

	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	if got["message"] == "" {
		t.Errorf("expected non-empty message in response")
	}
}
