package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	router.GET("/songs", getSongs)
	router.GET("/songs/:id", getSongByID)
	router.POST("/songs", addSong)

	router.POST("/quicksort", func(c *gin.Context) {
		var input struct {
			Array []int `json:"array"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
			return
		}
		sorted := sortArray(input.Array)
		c.IndentedJSON(http.StatusOK, gin.H{"sorted": sorted})
	})
}
