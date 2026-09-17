package api

import (
	"net/http"

	"github.com/Unleash/unleash-client-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/logging"
	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/service"
)

func authLoginHandler(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if !unleash.IsEnabled("login-feature-flag") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Funcionalidad no disponible"})
			return
		}

		var request model.LoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logging.FromContext(ctx).Warn("login_invalid_request", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := authService.Authenticate(ctx, request.Username, request.Password)
		if err != nil {
			// Nunca logueamos la password, solo el username y el motivo del fallo.
			logging.FromContext(ctx).Warn("login_failed", "username", request.Username, "error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		logging.FromContext(ctx).Info("login_succeeded", "username", request.Username)
		c.JSON(http.StatusOK, model.LoginResponse{Username: request.Username, Token: token})
	}
}

func flagStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		enabled := unleash.IsEnabled("login-feature-flag")
		c.JSON(http.StatusOK, gin.H{
			"feature_login": enabled,
			"feature_lab":   true,
		})
	}
}
