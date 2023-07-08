package dto

import (
	"time"
)

type (
	CreateVAResponse struct {
		ID             string    `json:"id"`
		OwnerID        string    `json:"owner_id"`
		ExternalID     string    `json:"external_id"`
		AccountNumber  string    `json:"account_number"`
		BankCode       string    `json:"bank_code"`
		MerchantCode   string    `json:"merchant_code"`
		Name           string    `json:"name"`
		IsClosed       bool      `json:"is_closed"`
		ExpectedAmount int64     `json:"expected_amount"`
		ExpirationDate time.Time `json:"expiration_date"`
		IsSingleUse    bool      `json:"is_single_use"`
		Status         string    `json:"status"`
		ErrorCode      string    `json:"error_code,omitempty"`
		ErrorMessage   string    `json:"message,omitempty"`
	}
)
