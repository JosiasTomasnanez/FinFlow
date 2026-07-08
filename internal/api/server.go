package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/service"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func NewServer(walletService *service.WalletService, authService *service.AuthService) *http.Server {
	router := gin.Default()

	router.Use(cors.Default())

	p := ginprometheus.NewPrometheus("gin")

	p.MetricsPath = "metrics"

	p.Use(router)

	apiGroup := router.Group("/api")
	apiGroup.GET("/health", healthHandler)
	apiGroup.GET("/wallets", walletListHandler(walletService))
	apiGroup.POST("/wallets", walletCreateHandler(walletService))
	apiGroup.GET("/wallets/:walletID", walletDetailHandler(walletService))
	apiGroup.POST("/payments", paymentHandler(walletService))
	apiGroup.POST("/login", authLoginHandler(authService))
	apiGroup.GET("/flags", flagStatusHandler())

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
