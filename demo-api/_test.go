package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// helper to create a test router with the same routes as main
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/songs", getSongs)
	return router
}

func TestGetSongs(t *testing.T) {
	// ensure gin is running in test mode so logs are suppressed
	gin.SetMode(gin.TestMode)

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
