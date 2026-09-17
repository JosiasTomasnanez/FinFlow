package api

import (
	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/lab"
)

func registerLabRoutes(apiGroup *gin.RouterGroup, mem *lab.Hold) {
	g := apiGroup.Group("/lab")
	g.GET("/ok", labOKHandler)
	g.POST("/slow", labSlowHandler)
	g.POST("/error", labErrorHandler)
	g.POST("/cpu", labCPUHandler)
	g.POST("/burst", labBurstHandler)
	g.POST("/memory", labMemoryRetainHandler(mem))
	g.POST("/memory/release", labMemoryReleaseHandler(mem))
}
