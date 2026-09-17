package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/lab"
)

type labSlowRequest struct {
	DelayMS int `json:"delay_ms"`
}

// labSlowHandler duerme y responde 200. El p95 de Grafana solo cuenta [23]xx.
func labSlowHandler(c *gin.Context) {
	var req labSlowRequest
	_ = c.ShouldBindJSON(&req)
	d := lab.DurationFromMS(req.DelayMS, lab.DefaultDelay, lab.MaxDelay)
	time.Sleep(d)
	c.JSON(http.StatusOK, gin.H{"slept_ms": d.Milliseconds()})
}
