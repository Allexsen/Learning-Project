package router

import (
	"github.com/Allexsen/Learning-Project/internal/handlers"
)

type Record struct {
	// Will get an UID later on
	Name          string
	Email         string
	HoursWorked   int
	MinutesWorked int
}

func initRecordsRouter() {
	r.POST("/record/add", handlers.RecordAdd())
}
