package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Attaching routes with REST verbs
	// gin.Context hold the information of the individual request. Can serialize
	// data into JSON before sending it back to the client using the context.JSON function
	r.GET("/pingTime", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"serverTime": time.Now().UTC(),
		})
	})
	r.Run(":8000") // Default listen and serve on 0.0.0.0:8000
}
