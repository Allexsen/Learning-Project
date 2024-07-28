package router

import (
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/gin-gonic/gin"
)

func initWsRouter(wsManager *ws.WsManager) {
	wsRouter := r.Group("/ws")
	{
		// TODO: Swap out the placeholder
		wsRouter.GET("", func(c *gin.Context) {
			ws.ServeWs(wsManager, c)
		})
	}
}
