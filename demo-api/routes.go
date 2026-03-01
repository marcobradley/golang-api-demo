package main

import "github.com/gin-gonic/gin"

func registerRoutes(router *gin.Engine) {
	router.GET("/songs", getSongs)
	router.GET("/songs/:id", getSongByID)
	router.POST("/songs", addSong)
}
