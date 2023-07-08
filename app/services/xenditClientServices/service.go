package xendit_client

import (
	"encoding/json"
	"ppob-backend/app/dto"
	"ppob-backend/config"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type (
	xenditClient struct {
		Config *config.SystemConfig
	}

	XenditClient interface {
		CreateVA(bank_code string, va_name string, amount int) (dto.CreateVAResponse, error)
	}

	xenditResponse struct {
		Data dto.CreateVAResponse `json:"data"`
	}

	createVAPayload struct {
		ExternalID     string    `json:"external_id"`
		BankCode       string    `json:"bank_code"`
		Name           string    `json:"name"`
		IsSingleUse    bool      `json:"is_single_use"`
		IsClosed       bool      `json:"is_closed"`
		ExpectedAmount int64     `json:"expected_amount"`
		ExpirationDate time.Time `json:"expiration_date"`
	}
)

func NewXenditClient(config *config.SystemConfig) XenditClient {
	return &xenditClient{
		Config: config,
	}
}

func (x *xenditClient) CreateVA(bank_code string, va_name string, amount int) (dto.CreateVAResponse, error) {
	var (
		body     = xenditResponse{}
		response = dto.CreateVAResponse{}
		client   = resty.New()
		now      = time.Now()
	)

	payload := createVAPayload{
		ExternalID:     uuid.NewString(),
		BankCode:       bank_code,
		Name:           va_name,
		IsSingleUse:    true,
		IsClosed:       true,
		ExpectedAmount: int64(amount),
		ExpirationDate: now.AddDate(0, 0, x.Config.XenditVAExpiration),
	}

	resp, err := client.R().SetBody(payload).Post(x.Config.XenditBaseUrl + x.Config.XenditCreateVAPath)
	if err != nil {
		x.Config.Logger.Warnf("failed to call API: %v", err)
		return response, err
	}

	_ = json.Unmarshal(resp.Body(), &body)

	response = dto.CreateVAResponse(body.Data)

	return response, nil
}

// (x *xenditClient) CreateVA(bank_code string, va_name string, amount int) (dto.CreateVAResponse, error) {
