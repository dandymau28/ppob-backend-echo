package dto

import "time"

type (
	GetBalanceRequest struct {
		UserID string `param:"user_id"`
	}

	GetBalanceResponse struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Balance  int64  `json:"balance"`
	}

	TopupWalletRequest struct {
		UserID string `param:"user_id"`
		Amount int64  `json:"amount"`
	}

	TopupWalletResponse struct {
		Name          string `json:"name"`
		Username      string `json:"username"`
		BalanceBefore int64  `json:"balance_before"`
		Balance       int64  `json:"balance"`
	}

	GetUserBalance struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Balance  int64  `json:"balance"`
	}

	DoTransactionResponse struct {
		Balance       int64     `json:"balance"`
		Status        string    `json:"status"`
		TransactionAt time.Time `json:"transaction_at"`
		ProductName   string    `json:"product_name"`
		Price         int64     `json:"price"`
	}

	BuyPulsaRequest struct {
		UserID      string `json:"user_id"`
		ProductCode string `json:"product_code"`
		CustomerNo  string `json:"customer_no"`
	}

	TransactionHistoryRequest struct {
		UserID string `param:"user_id"`
	}

	Product struct {
		ProductName string `json:"product_name"`
		ProductCode string `json:"product_code"`
		Price       int64  `json:"price"`
	}

	PrePurchaseResponse struct {
		TransactionID string  `json:"transaction_id"`
		Product       Product `json:"product"`
		Total         int64   `json:"total"`
	}
)
