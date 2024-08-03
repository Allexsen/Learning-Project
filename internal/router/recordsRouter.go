// Package router sets up API endpoints to serve requests
package router

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/handlers"
	"github.com/Allexsen/Learning-Project/internal/middlewares"
)

// initRecordsRouter sets up routes associated with records
func initRecordsRouter() {
	log.Println("Setting up records router...")
	recordRouter := r.Group("/record")
	recordRouter.Use(middlewares.CheckJWT())
	{
		recordRouter.POST("/add", handlers.RecordAdd)
		recordRouter.POST("/delete", handlers.RecordDelete)
	}
}
