package model

// PaymentRequest defines a transfer request between wallets.
type PaymentRequest struct {
    FromWalletID string `json:"from_wallet_id"`
    ToWalletID   string `json:"to_wallet_id"`
    Amount       int64  `json:"amount"`
}

// PaymentResult returns the result of a wallet transfer.
type PaymentResult struct {
    FromWalletID string `json:"from_wallet_id"`
    ToWalletID   string `json:"to_wallet_id"`
    Amount       int64  `json:"amount"`
    Status       string `json:"status"`
}
