package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	registerRoutes(router)

	// bind to all addresses so the service is reachable from outside the container
	router.Run(":8080")
}
