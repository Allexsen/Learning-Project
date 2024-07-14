// Package router sets up API endpoints to serve requests
package router

import (
	"github.com/gin-gonic/gin"
)

var (
	// Gin engine used through the application
	r = gin.Default()
)

// GetEngine returns the router of the application
func GetEngine() *gin.Engine {
	return r
}

// InitRouter initializes default rout and statics folder routing,
// and then invokes initialization of other routers too.
func InitRouter() {
	r.Static("/statics/", "../../static/")

	r.GET("/", func(c *gin.Context) {
		c.File("../../static/html/index.html")
	})

	initUserRouter()
	initRecordsRouter()

	r.Run(":8080")
}
