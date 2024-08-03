package router

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/middlewares"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/gin-gonic/gin"
)

func initWsRouter(wsManager *ws.WsManager) {
	log.Println("Setting up ws router...")
	wsRouter := r.Group("/ws")
	wsRouter.Use(middlewares.CheckJWT())
	{
		// TODO: Swap out the placeholder
		wsRouter.GET("", func(c *gin.Context) {
			ws.WsHandler(wsManager, c)
		})
	}
}
