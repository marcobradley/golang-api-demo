package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type song struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var songs = []song{
	{ID: "1", Title: "Shape of You", Artist: "Ed Sheeran", Price: 1.29},
	{ID: "2", Title: "Blinding Lights", Artist: "The Weeknd", Price: 1.29},
	{ID: "3", Title: "Dance Monkey", Artist: "Tones and I", Price: 1.29},
}

func getSongs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, songs)
}

func getSongByID(c *gin.Context) {
	id := c.Param("id")
	for _, s := range songs {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "song not found"})
}

func addSong(c *gin.Context) {
	var newSong song
	if err := c.BindJSON(&newSong); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	songs = append(songs, newSong)
	c.IndentedJSON(http.StatusCreated, newSong)
}
