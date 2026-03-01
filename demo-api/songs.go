package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type song struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var initialSongs = []song{
	{ID: "1", Title: "Shape of You", Artist: "Ed Sheeran", Price: 1.29},
	{ID: "2", Title: "Blinding Lights", Artist: "The Weeknd", Price: 1.29},
	{ID: "3", Title: "Dance Monkey", Artist: "Tones and I", Price: 1.29},
}

var (
	mu    sync.RWMutex
	songs = append([]song{}, initialSongs...)
)

func getSongs(c *gin.Context) {
	mu.RLock()
	snapshot := make([]song, len(songs))
	copy(snapshot, songs)
	mu.RUnlock()
	c.IndentedJSON(http.StatusOK, snapshot)
}

func getSongByID(c *gin.Context) {
	id := c.Param("id")
	mu.RLock()
	var found *song
	for i := range songs {
		if songs[i].ID == id {
			foundSong := songs[i]
			found = &foundSong
			break
		}
	}
	mu.RUnlock()
	if found != nil {
		c.IndentedJSON(http.StatusOK, *found)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "song not found"})
}

func addSong(c *gin.Context) {
	var newSong song
	if err := c.BindJSON(&newSong); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	if newSong.ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is required"})
		return
	}
	mu.Lock()
	for _, s := range songs {
		if s.ID == newSong.ID {
			mu.Unlock()
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "song with this id already exists"})
			return
		}
	}
	songs = append(songs, newSong)
	mu.Unlock()
	c.IndentedJSON(http.StatusCreated, newSong)
}
