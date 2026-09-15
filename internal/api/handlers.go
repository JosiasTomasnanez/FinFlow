package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josiastomasnanez/finflow/internal/logging"
	"github.com/josiastomasnanez/finflow/internal/model"
	"github.com/josiastomasnanez/finflow/internal/service"
)

func walletListHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		wallets := service.ListWallets(c.Request.Context())
		c.JSON(http.StatusOK, wallets)
	}
}

func walletCreateHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var request model.WalletCreateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logging.FromContext(ctx).Warn("wallet_create_invalid_request", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		wallet, err := service.CreateWallet(ctx, request.Owner, request.InitialBalance)
		if err != nil {
			logging.FromContext(ctx).Warn("wallet_create_failed", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logging.FromContext(ctx).Info("wallet_created", "wallet_id", wallet.ID, "owner", wallet.Owner)
		c.JSON(http.StatusCreated, wallet)
	}
}

func walletDetailHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		walletID := c.Param("walletID")
		wallet, found := service.GetWallet(ctx, walletID)
		if !found {
			logging.FromContext(ctx).Warn("wallet_not_found", "wallet_id", walletID)
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}

		c.JSON(http.StatusOK, wallet)
	}
}

func paymentHandler(service *service.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var request model.PaymentRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logging.FromContext(ctx).Warn("payment_invalid_request", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment, err := service.Transfer(ctx, request.FromWalletID, request.ToWalletID, request.Amount)
		if err != nil {
			logging.FromContext(ctx).Warn("payment_failed",
				"error", err.Error(),
				"from_wallet_id", request.FromWalletID,
				"to_wallet_id", request.ToWalletID,
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logging.FromContext(ctx).Info("payment_completed",
			"from_wallet_id", payment.FromWalletID,
			"to_wallet_id", payment.ToWalletID,
			"amount", payment.Amount,
		)
		c.JSON(http.StatusCreated, payment)
	}
}
