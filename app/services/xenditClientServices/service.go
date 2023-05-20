package xendit_client

import (
	"ppob-backend/app/dto"
	"ppob-backend/config"

	"github.com/labstack/echo/v4"
)

type (
	xenditClient struct {
		Config  *config.SystemConfig
		context echo.Context
	}

	XenditClient interface {
		Topup(payload dto.BuyProductPayload) (dto.BuyProductResponse, error)
	}

	xenditResponse struct {
		Data dto.BuyProductResponse `json:"data"`
	}
)

func NewDfWebClient(config *config.SystemConfig) DfWebClient {
	return &dfWebClient{
		Config: config,
	}
}
