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
	return []song{
		{ID: "1", Title: "Shape of You", Artist: "Ed Sheeran", Price: 1.29},
		{ID: "2", Title: "Blinding Lights", Artist: "The Weeknd", Price: 1.29},
		{ID: "3", Title: "Dance Monkey", Artist: "Tones and I", Price: 1.29},
	}
}

func resetSongsForTest(t *testing.T) {
	t.Helper()
	songs = seedSongs()
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
	if len(got) != len(songs) {
		t.Errorf("expected %d songs but got %d", len(songs), len(got))
	}

	// simple field check
	for i := range songs {
		if got[i] != songs[i] {
			t.Errorf("song at index %d does not match\ngot: %#v\nwant: %#v", i, got[i], songs[i])
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

	if got != songs[0] {
		t.Errorf("returned song does not match\ngot: %#v\nwant: %#v", got, songs[0])
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

	if len(songs) != 4 {
		t.Fatalf("expected 4 songs after add but got %d", len(songs))
	}

	last := songs[len(songs)-1]
	if last != got {
		t.Errorf("last song in collection does not match created song\ngot: %#v\nwant: %#v", last, got)
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
	originalCount := len(songs)

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

	if len(songs) != originalCount {
		t.Errorf("songs collection changed on invalid request: got %d, want %d", len(songs), originalCount)
	}
}
