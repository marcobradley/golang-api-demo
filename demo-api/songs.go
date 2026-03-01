package main

import (
	"net/http"
	"sort"
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
	idx := findSongIndexByID(songs, id)
	var foundSong song
	found := idx < len(songs) && songs[idx].ID == id
	if found {
		foundSong = songs[idx]
	}
	mu.RUnlock()
	if found {
		c.IndentedJSON(http.StatusOK, foundSong)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "song not found"})
}

func findSongIndexByID(list []song, id string) int {
	return sort.Search(len(list), func(i int) bool {
		return list[i].ID >= id
	})
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
	idx := findSongIndexByID(songs, newSong.ID)
	if idx < len(songs) && songs[idx].ID == newSong.ID {
		mu.Unlock()
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "song with this id already exists"})
		return
	}
	songs = append(songs, song{})
	copy(songs[idx+1:], songs[idx:])
	songs[idx] = newSong
	mu.Unlock()
	c.IndentedJSON(http.StatusCreated, newSong)
}
