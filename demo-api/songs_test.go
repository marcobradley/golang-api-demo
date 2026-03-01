package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func seedSongs() []song {
	cloned := make([]song, len(initialSongs))
	copy(cloned, initialSongs)
	return cloned
}

func resetSongsForTest(t *testing.T) {
	t.Helper()
	mu.Lock()
	defer mu.Unlock()
	songs = seedSongs()
}

func songsCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(songs)
}

func songsSnapshot() []song {
	mu.RLock()
	defer mu.RUnlock()
	snapshot := make([]song, len(songs))
	copy(snapshot, songs)
	return snapshot
}

// helper to create a test router with the same routes as main
func setupRouter() *gin.Engine {
	router := gin.Default()
	registerRoutes(router)
	return router
}

func TestGetSongs(t *testing.T) {
	// ensure gin is running in test mode so logs are suppressed
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()

	req, err := http.NewRequest(http.MethodGet, "/songs", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, w.Code)
	}

	var got []song
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	// compare length first
	snapshot := songsSnapshot()
	if len(got) != len(snapshot) {
		t.Errorf("expected %d songs but got %d", len(snapshot), len(got))
	}

	// simple field check
	for i := range snapshot {
		if got[i] != snapshot[i] {
			t.Errorf("song at index %d does not match\ngot: %#v\nwant: %#v", i, got[i], snapshot[i])
		}
	}
}

func TestGetSongByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()

	req, err := http.NewRequest(http.MethodGet, "/songs/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, w.Code)
	}

	var got song
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	snapshot := songsSnapshot()
	if got != snapshot[0] {
		t.Errorf("returned song does not match\ngot: %#v\nwant: %#v", got, snapshot[0])
	}
}

func TestAddSong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()

	body := []byte(`{"id":"4","title":"Levitating","artist":"Dua Lipa","price":1.49}`)
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d but got %d", http.StatusCreated, w.Code)
	}

	var got song
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	if got.ID != "4" || got.Title != "Levitating" || got.Artist != "Dua Lipa" || got.Price != 1.49 {
		t.Errorf("created song payload mismatch: %#v", got)
	}

	if songsCount() != 4 {
		t.Fatalf("expected 4 songs after add but got %d", songsCount())
	}

	snapshot := songsSnapshot()
	if snapshot[3] != got {
		t.Errorf("song at index %d does not match created song\ngot: %#v\nwant: %#v", 3, snapshot[3], got)
	}
}

func TestAddSongMaintainsAscendingIDOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()

	body := []byte(`{"id":"0","title":"Before All","artist":"Tester","price":0.49}`)
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d but got %d", http.StatusCreated, w.Code)
	}

	if songsCount() != 4 {
		t.Fatalf("expected 4 songs after add but got %d", songsCount())
	}

	snapshot := songsSnapshot()
	if snapshot[0].ID != "0" {
		t.Fatalf("expected first song ID %q but got %q", "0", snapshot[0].ID)
	}

	for i := 1; i < len(snapshot); i++ {
		if snapshot[i-1].ID > snapshot[i].ID {
			t.Fatalf("songs are not in ascending ID order at index %d: %q > %q", i, snapshot[i-1].ID, snapshot[i].ID)
		}
	}
}

func TestGetSongByIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()

	req, err := http.NewRequest(http.MethodGet, "/songs/999", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d but got %d", http.StatusNotFound, w.Code)
	}

	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	if got["message"] != "song not found" {
		t.Errorf("expected message %q but got %q", "song not found", got["message"])
	}
}

func TestAddSongInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()
	originalCount := songsCount()

	body := []byte(`{"id":`)
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(body))
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

	if got["message"] != "invalid request body" {
		t.Errorf("expected message %q but got %q", "invalid request body", got["message"])
	}

	if songsCount() != originalCount {
		t.Errorf("songs collection changed on invalid request: got %d, want %d", songsCount(), originalCount)
	}
}

func TestAddSongMissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()
	originalCount := songsCount()

	body := []byte(`{"title":"No ID Song","artist":"Unknown","price":0.99}`)
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(body))
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

	if got["message"] != "id is required" {
		t.Errorf("expected message %q but got %q", "id is required", got["message"])
	}

	if songsCount() != originalCount {
		t.Errorf("songs collection changed on missing id: got %d, want %d", songsCount(), originalCount)
	}
}

func TestAddSongDuplicateID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resetSongsForTest(t)

	router := setupRouter()
	originalCount := songsCount()

	body := []byte(`{"id":"1","title":"Duplicate","artist":"Someone","price":0.99}`)
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d but got %d", http.StatusConflict, w.Code)
	}

	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body not valid json: %v", err)
	}

	if got["message"] != "song with this id already exists" {
		t.Errorf("expected message %q but got %q", "song with this id already exists", got["message"])
	}

	if songsCount() != originalCount {
		t.Errorf("songs collection changed on duplicate id: got %d, want %d", songsCount(), originalCount)
	}
}
