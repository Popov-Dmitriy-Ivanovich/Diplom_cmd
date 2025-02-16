package routes

import (
	"github.com/gin-gonic/gin"
)

func WriteRoutes(rg *gin.RouterGroup, rw ...RouteWriter) {
	for _, rw := range rw {
		rw.WriteRoutes(rg)
	}
}

type RouteWriter interface {
	WriteRoutes(*gin.RouterGroup)
}