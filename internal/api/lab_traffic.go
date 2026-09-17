package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// labOKHandler: 200 inmediato. El front dispara N en paralelo para RPS.
func labOKHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
