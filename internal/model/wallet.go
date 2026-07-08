package model

type Wallet struct {
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

type WalletCreateRequest struct {
	Owner          string `json:"owner"`
	InitialBalance int64  `json:"initial_balance"`
}
