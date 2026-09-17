package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/lab"
)

type labCPURequest struct {
	DurationMS int `json:"duration_ms"`
}

type labBurstRequest struct {
	SpinMS int `json:"spin_ms"`
}

type labMemoryRequest struct {
	Megabytes int `json:"megabytes"`
}

func labCPUHandler(c *gin.Context) {
	var req labCPURequest
	_ = c.ShouldBindJSON(&req)
	d := lab.DurationFromMS(req.DurationMS, lab.DefaultCPU, lab.MaxCPU)
	lab.Spin(d)
	c.JSON(http.StatusOK, gin.H{"spun_ms": d.Milliseconds()})
}

func labBurstHandler(c *gin.Context) {
	var req labBurstRequest
	_ = c.ShouldBindJSON(&req)
	d := lab.DurationFromMS(req.SpinMS, lab.DefaultBurst, lab.MaxBurst)
	lab.Spin(d)
	c.JSON(http.StatusOK, gin.H{"spun_ms": d.Milliseconds()})
}

func labMemoryRetainHandler(mem *lab.Hold) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req labMemoryRequest
		_ = c.ShouldBindJSON(&req)
		held := mem.Retain(req.Megabytes)
		c.JSON(http.StatusOK, gin.H{"held_mib": held})
	}
}

func labMemoryReleaseHandler(mem *lab.Hold) gin.HandlerFunc {
	return func(c *gin.Context) {
		mem.Release()
		c.JSON(http.StatusOK, gin.H{"held_mib": 0})
	}
}
