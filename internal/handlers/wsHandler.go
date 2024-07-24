package handlers

import (
	"github.com/Allexsen/Learning-Project/internal/ws"
	"github.com/gin-gonic/gin"
)

func ServeWs(c *gin.Context) {
	wsManager := ws.NewWsManager()
	go wsManager.Run()

	ws.ServeWs(wsManager, c)
}
