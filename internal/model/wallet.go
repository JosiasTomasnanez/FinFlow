package model

// Wallet represents a basic fintech wallet.
type Wallet struct {
    ID      string `json:"id"`
    Owner   string `json:"owner"`
    Balance int64  `json:"balance"`
}

// WalletCreateRequest contains the payload to create a new wallet.
type WalletCreateRequest struct {
    Owner          string `json:"owner"`
    InitialBalance int64  `json:"initial_balance"`
}
