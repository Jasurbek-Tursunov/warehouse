package pgin

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"client": c.ClientIP(), "server": c.Request.Host})
	})
	return router
}
