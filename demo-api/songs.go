package main

import (
	"math"
	"net/http"
	"sort"
	"strconv"
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
	targetID, err := strconv.Atoi(id)
	if err != nil {
		return len(list)
	}

	return sort.Search(len(list), func(i int) bool {
		return numericSongIDValue(list[i].ID) >= targetID
	})
}

func numericSongIDValue(id string) int {
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return math.MaxInt
	}
	return parsedID
}

func nextSongID(list []song) string {
	maxID := 0
	for _, currentSong := range list {
		parsedID, err := strconv.Atoi(currentSong.ID)
		if err != nil {
			continue
		}
		if parsedID > maxID {
			maxID = parsedID
		}
	}
	return strconv.Itoa(maxID + 1)
}

func addSong(c *gin.Context) {
	var newSong song
	if err := c.ShouldBindJSON(&newSong); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	mu.Lock()
	newSong.ID = nextSongID(songs)
	idx := findSongIndexByID(songs, newSong.ID)
	songs = append(songs, song{})
	copy(songs[idx+1:], songs[idx:])
	songs[idx] = newSong
	mu.Unlock()
	c.IndentedJSON(http.StatusCreated, newSong)
}
