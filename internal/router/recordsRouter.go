package router

import (
	"github.com/Allexsen/Learning-Project/internal/handlers"
)

func initRecordsRouter() {
	r.POST("/record/add", handlers.RecordAdd())
	r.POST("/record/delete", handlers.RecordDelete())
}
