package router

import (
	"github.com/gin-gonic/gin"
)

var (
	r = gin.Default()
)

func GetEngine() *gin.Engine {
	return r
}

func InitRouter() {
	r.Static("/statics", "../../static/")

	r.GET("/", func(c *gin.Context) {
		c.File("../../static/index.html")
	})

	initUserRouter()
	initRecordsRouter()

	r.Run(":8080")
}
