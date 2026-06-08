package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// 1. Agregamos el import de Prometheus para Gin
	"github.com/josiastomasnanez/finflow/internal/service"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// NewServer returns a configured HTTP server for the FinFlow API.
func NewServer(walletService *service.WalletService, authService *service.AuthService) *http.Server {
	router := gin.Default()

	router.Use(cors.Default())

	// 2. CONFIGURACIÓN DE PROMETHEUS
	p := ginprometheus.NewPrometheus("gin")

	// Le indicamos que exponga las métricas en la ruta "/api/metrics"
	// para que coincida con tu apiGroup y Nginx lo capture bien
	p.MetricsPath = "metrics"

	// 3. Vinculamos Prometheus al router de Gin
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
