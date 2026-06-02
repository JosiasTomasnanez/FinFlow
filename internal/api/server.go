package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/josiastomasnanez/finflow/internal/service"
)

// NewServer returns a configured HTTP server for the FinFlow API.
func NewServer(service *service.WalletService) *http.Server {
    router := gin.Default()

    apiGroup := router.Group("/api")
    apiGroup.GET("/health", healthHandler)
    apiGroup.GET("/wallets", walletListHandler(service))
    apiGroup.POST("/wallets", walletCreateHandler(service))
    apiGroup.GET("/wallets/:walletID", walletDetailHandler(service))
    apiGroup.POST("/payments", paymentHandler(service))

    router.Static("/assets", "./frontend/dist/assets")
    router.StaticFile("/", "./frontend/dist/index.html")
    router.NoRoute(func(c *gin.Context) {
        c.File("./frontend/dist/index.html")
    })

    return &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
}

func healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "FinFlow"})
}
