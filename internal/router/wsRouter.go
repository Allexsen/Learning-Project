package router

import (
	"github.com/Allexsen/Learning-Project/internal/handlers"
)

func initWsRouter() {
	wsRouter := r.Group("/ws")
	{
		// Placeholder for now
		wsRouter.GET("", handlers.ServeWs)
	}
}
