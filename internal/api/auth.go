package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/service"
)

func authLoginHandler(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.LoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := authService.Authenticate(request.Username, request.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.LoginResponse{Username: request.Username, Token: token})
	}
}

func flagStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		enabled := os.Getenv("FEATURE_LOGIN") == "true"
		c.JSON(http.StatusOK, gin.H{"feature_login": enabled})
	}
}
