package api

import (
	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/lab"
)

type labErrorRequest struct {
	Code int `json:"code"`
}

func labErrorHandler(c *gin.Context) {
	var req labErrorRequest
	_ = c.ShouldBindJSON(&req)
	status := lab.AllowedErrorStatus(req.Code)
	c.JSON(status, gin.H{"error": "lab_injected_error", "code": status})
}
