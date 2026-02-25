package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/songs", getSongs)

	// bind to all addresses so the service is reachable from outside the container
	router.Run(":8080")
}

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
