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
	r.Static("/static", "../../static/")

	r.GET("/", func(c *gin.Context) {
		c.File("./index.html")
	})

	initRecordsRouter()

}
