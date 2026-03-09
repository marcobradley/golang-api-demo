package main

import "github.com/gin-gonic/gin"

func registerRoutes(router *gin.Engine) {
	router.GET("/songs", getSongs)
	router.GET("/songs/:id", getSongByID)
	router.POST("/songs", addSong)

	router.POST("/quicksort", func(c *gin.Context) {
		var input struct {
			Array []int `json:"array"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		sorted := sortArray(input.Array)
		c.JSON(200, gin.H{"sorted": sorted})
	})
}
