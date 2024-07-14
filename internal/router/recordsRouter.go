// Package router sets up API endpoints to serve requests
package router

import (
	"github.com/Allexsen/Learning-Project/internal/handlers"
	"github.com/Allexsen/Learning-Project/internal/middlewares"
)

// initRecordsRouter sets up routes associated with records
func initRecordsRouter() {
	records := r.Group("/record")
	records.Use(middlewares.CheckJWT())
	{
		r.POST("/record/add", handlers.RecordAdd)
		r.POST("/record/delete", handlers.RecordDelete)
	}
}
