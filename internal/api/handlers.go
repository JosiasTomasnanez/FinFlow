package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/service"
)

func walletListHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wallets := service.ListWallets()
		c.JSON(http.StatusOK, wallets)
	}
}

func walletCreateHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.WalletCreateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		wallet, err := service.CreateWallet(request.Owner, request.InitialBalance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, wallet)
	}
}

func walletDetailHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletID")
		wallet, found := service.GetWallet(walletID)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}

		c.JSON(http.StatusOK, wallet)
	}
}

func paymentHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request model.PaymentRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment, err := service.Transfer(request.FromWalletID, request.ToWalletID, request.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, payment)
	}
}
