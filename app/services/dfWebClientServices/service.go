package dfwebclientservices

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"ppob-backend/app/dto"
	"ppob-backend/config"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type (
	dfWebClient struct {
		Config  *config.SystemConfig
		context echo.Context
	}

	DfWebClient interface {
		Topup(payload dto.BuyProductPayload) (dto.BuyProductResponse, error)
	}

	dfResponse struct {
		Data dto.BuyProductResponse `json:"data"`
	}
)

func NewDfWebClient(config *config.SystemConfig) DfWebClient {
	return &dfWebClient{
		Config: config,
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetRefID(code string) string {
	numberCode := rand.Intn(999999)
	trxTime := time.Now().UnixMicro()
	refID := fmt.Sprintf("%s%d%06d", code, trxTime, numberCode)
	return refID
}

func (s *dfWebClient) Topup(payload dto.BuyProductPayload) (dto.BuyProductResponse, error) {
	var (
		body     = dfResponse{}
		response = dto.BuyProductResponse{}
		client   = resty.New()
	)

	payload.Username = s.Config.DigiflazzUsername
	payload.RefId = GetRefID("GPOB")
	payload.Sign = GetMD5Hash(s.Config.DigiflazzUsername + s.Config.DigiflazzApiKey + payload.RefId)
	payload.Testing = false

	if s.Config.DigiflazzTesting == "1" {
		payload.Testing = true
	}

	resp, err := client.R().SetBody(payload).Post(s.Config.DigiflazzBaseUrl + s.Config.DigiflazzTopupPath)
	if err != nil {
		s.context.Logger().Warnf("failed to call API: %v", err)
		return response, err
	}

	_ = json.Unmarshal(resp.Body(), &body)

	response = dto.BuyProductResponse(body.Data)

	return response, nil
}
