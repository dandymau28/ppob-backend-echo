package dto

type (
	BuyProductPayload struct {
		Username     string `json:"username"`
		BuyerSkuCode string `json:"buyer_sku_code"`
		CustomerNo   string `json:"customer_no"`
		RefId        string `json:"ref_id"`
		Sign         string `json:"sign"`
		Testing      bool   `json:"testing"`
	}

	BuyProductResponse struct {
		RefId          string  `json:"ref_id"`
		CustomerNo     string  `json:"customer_no"`
		BuyerSkuCode   string  `json:"buyer_sku_code"`
		Message        string  `json:"message"`
		Status         string  `json:"status"`
		Rc             string  `json:"rc"`
		Sn             string  `json:"sn"`
		BuyerLastSaldo float64 `json:"buyer_last_saldo"`
		Price          int64   `json:"price"`
	}

	WebhookHeaders struct {
		XDigiflazzDelivery string `header:"X-Digiflazz-Delivery"`
		XHubSignature      string `header:"X-Hub-Signature"`
		XDigiflazzEvent    string `header:"X-Digiflazz-Event"`
	}

	WebhookRequest struct {
		Data WebhookRequestBody `json:"data"`
	}

	WebhookRequestBody struct {
		TrxID          string `json:"trx_id"`
		RefID          string `json:"ref_id"`
		CustomerNo     string `json:"customer_no"`
		BuyerSkuCode   string `json:"buyer_sku_code"`
		Message        string `json:"message"`
		Status         string `json:"status"`
		Rc             string `json:"rc"`
		BuyerLastSaldo int64  `json:"buyer_last_saldo"`
		Sn             string `json:"sn"`
		Price          int64  `json:"price"`
		Tele           string `json:"tele"`
		Wa             string `json:"wa"`
	}
)
