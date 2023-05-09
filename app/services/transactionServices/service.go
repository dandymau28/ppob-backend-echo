package transactionServices

import (
	"errors"
	"fmt"
	"ppob-backend/app/dto"
	"ppob-backend/app/repository/transactionRepository"
	dfwebclientservices "ppob-backend/app/services/dfWebClientServices"
	"ppob-backend/config"
	"ppob-backend/model"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	transactionServices struct {
		txnRepository transactionRepository.TransactionRepository
		Config        *config.SystemConfig
		ctx           echo.Context
		Mutex         sync.Mutex
		dfWebClient   dfwebclientservices.DfWebClient
		TAG           string
	}

	TransactionServices interface {
		GetUserBalanceByUserID(request dto.GetBalanceRequest) (dto.GetBalanceResponse, error)
		TopupWallet(request dto.TopupWalletRequest) (dto.TopupWalletResponse, error)
		PrePurchase(product_code string, customer_no string, user_id string) (dto.PrePurchaseResponse, error)
		DoTransaction(ref_id string, user_id string) (dto.DoTransactionResponse, error)
		GetTransactionHistory(user_id string) ([]model.Transaction, error)
		WebhookHandle(headers dto.WebhookHeaders, body dto.WebhookRequestBody) error
	}
)

func NewTransactionServices(config *config.SystemConfig, txnRepo transactionRepository.TransactionRepository, dfWebClient dfwebclientservices.DfWebClient) TransactionServices {
	return &transactionServices{
		txnRepository: txnRepo,
		Config:        config,
		dfWebClient:   dfWebClient,
		TAG:           "services.transaction-services",
		ctx:           echo.New().AcquireContext(),
	}
}

func (s *transactionServices) GetUserBalanceByUserID(request dto.GetBalanceRequest) (dto.GetBalanceResponse, error) {
	var (
		response = dto.GetBalanceResponse{}
	)

	userID := request.UserID

	userBalance := s.txnRepository.GetUserBalance(userID)

	if (dto.GetUserBalance{}) == userBalance {
		return response, errors.New("no user balance found")
	}

	response = dto.GetBalanceResponse(userBalance)

	return response, nil
}

func (s *transactionServices) TopupWallet(request dto.TopupWalletRequest) (dto.TopupWalletResponse, error) {
	defer echo.New().ReleaseContext(s.ctx)

	s.ctx.Logger().Infof("%s Start Topup Wallet for %s", s.TAG, request.UserID)
	var (
		response = dto.TopupWalletResponse{}
		wallet   = model.Wallet{}
	)

	userID := request.UserID

	wallet.Uuid = userID

	s.Mutex.Lock()
	err := s.txnRepository.GetWalletByUserID(&wallet)
	balanceBefore := wallet.Balance

	if err != nil {
		s.Mutex.Unlock()
		s.ctx.Logger().Infof("%v err getting wallet %v number %v", s.TAG, err, request.UserID)
		return response, errors.New("failed to get wallet")
	}

	wallet.Balance += int64(request.Amount)

	if wallet.Balance < 0 {
		s.Mutex.Unlock()
		s.ctx.Logger().Infof("%v err number %v insufficient balance %v", s.TAG, request.UserID, wallet.Balance)
		return response, errors.New("insufficient balance")
	}

	err = s.txnRepository.UpdateWallet(&wallet)

	if err != nil {
		s.Mutex.Unlock()
		s.ctx.Logger().Infof("%v err number %v update wallet %v", s.TAG, request.UserID, err)
		return response, errors.New("failed to update wallet")
	}
	s.Mutex.Unlock()

	userBalance := s.txnRepository.GetUserBalance(userID)

	response = dto.TopupWalletResponse{
		Name:          userBalance.Name,
		Username:      userBalance.Username,
		Balance:       userBalance.Balance,
		BalanceBefore: balanceBefore,
	}

	return response, nil
}

func (s *transactionServices) DoTransaction(trx_id string, user_id string) (dto.DoTransactionResponse, error) {
	/*
		- get wallet
		- get simulated transaction
		- if wallet balance less than price, failed
		- sub wallet balance to price
		- call digiflazz api
		- return response
	*/
	var (
		response     = dto.DoTransactionResponse{}
		balanceAfter int64
	)

	userBalance := s.txnRepository.GetUserBalance(user_id)

	txn := s.txnRepository.GetTransactionByTrxID(trx_id)

	if txn.Status != "SIMULATED" {
		return response, errors.New("invalid transaction")
	}

	if txn.PricePaid > userBalance.Balance {
		return response, errors.New("insufficient wallet balance")
	}

	topupWallet := dto.TopupWalletRequest{
		UserID: fmt.Sprint(user_id),
		Amount: txn.PricePaid * -1,
	}

	topupResp, err := s.TopupWallet(topupWallet)

	if err != nil {
		return response, err
	}

	payload := dto.BuyProductPayload{
		BuyerSkuCode: txn.BuyerSkuCode,
		CustomerNo:   txn.CustomerNo,
	}

	respDf, err := s.dfWebClient.Topup(payload)

	balanceAfter = topupResp.Balance

	if err != nil {
		refundWallet := dto.TopupWalletRequest{
			UserID: fmt.Sprint(user_id),
			Amount: txn.PricePaid,
		}

		//TODO: need to be enhanced when refund failed
		s.TopupWallet(refundWallet)

		s.ctx.Logger().Warnf("failed on calling api: %v", err)
		return response, err
	}

	if respDf.Rc != "00" && respDf.Status == "Gagal" {
		refundWallet := dto.TopupWalletRequest{
			UserID: fmt.Sprint(user_id),
			Amount: txn.PricePaid,
		}

		respRefund, _ := s.TopupWallet(refundWallet)

		balanceAfter = respRefund.Balance
	}

	txn.BalanceBefore = topupResp.BalanceBefore
	txn.BalanceAfter = balanceAfter
	txn.RefId = respDf.RefId
	txn.ResponseCode = respDf.Rc
	txn.ResponseMessage = respDf.Message
	txn.Sn = respDf.Sn
	txn.Status = strings.ToUpper(respDf.Status)
	txn.UpdatedAt = time.Now()

	err = s.txnRepository.UpdateTransaction(&txn)

	if err != nil {
		s.ctx.Logger().Warnf("failed to save transaction: %v", err)
		return response, err
	}

	response.Balance = txn.BalanceAfter
	response.Price = txn.PricePaid
	response.ProductName = txn.BuyerSkuCode
	response.Status = txn.Status
	response.TransactionAt = txn.CreatedAt

	return response, nil
}

func (s *transactionServices) GetTransactionHistory(user_id string) ([]model.Transaction, error) {
	var (
		response = []model.Transaction{}
	)

	response, err := s.txnRepository.GetTransactionHistoryByUserID(user_id)

	if err != nil {
		s.ctx.Logger().Warnf("failed to get transaction history: %v", err)
		return response, err
	}

	return response, nil
}

func (s *transactionServices) PrePurchase(product_code string, customer_no string, user_id string) (dto.PrePurchaseResponse, error) {
	/*
		- get wallet
		- get product price
		- if wallet balance less than price, failed
		- sub wallet balance to price
		- call digiflazz api
		- return response
	*/
	var (
		response = dto.PrePurchaseResponse{}
	)

	userBalance := s.txnRepository.GetUserBalance(user_id)

	product := s.txnRepository.GetProductByProductCode(product_code)

	if product.SellerPrice > userBalance.Balance {
		return response, errors.New("insufficient wallet balance")
	}

	txn := model.Transaction{}

	txnId := uuid.New()

	txn.BuyerSkuCode = product_code
	txn.CustomerNo = customer_no
	txn.PriceDist = product.Price
	txn.PricePaid = product.SellerPrice
	txn.UserId = user_id
	txn.CreatedAt = time.Now()
	txn.TransactionId = txnId.String()
	txn.Status = "SIMULATED"

	err := s.txnRepository.SaveTransaction(&txn)

	if err != nil {
		s.ctx.Logger().Warnf("failed to save transaction: %v", err)
		return response, err
	}

	response.Product.Price = product.SellerPrice
	response.Product.ProductCode = product.BuyerSkuCode
	response.Product.ProductName = product.ProductName
	response.Total = product.SellerPrice
	response.TransactionID = txn.TransactionId

	return response, nil
}

func (s *transactionServices) WebhookHandle(headers dto.WebhookHeaders, body dto.WebhookRequestBody) error {

	/*

		- validate hmac
		- !check header df event
		- get transaction by ref id
		- update transaction as in webhook body

	*/

	txn := s.txnRepository.GetTransactionByRefID(body.RefID)

	if txn.CustomerNo == "" {
		s.ctx.Logger().Warnf("no transaction found")
		return nil
	}

	txn.Status = strings.ToUpper(body.Status)
	txn.Sn = body.Sn
	txn.ResponseCode = body.Rc
	txn.ResponseMessage = body.Message

	if txn.Status == "Gagal" {
		refundWallet := dto.TopupWalletRequest{
			UserID: fmt.Sprint(txn.UserId),
			Amount: txn.PricePaid,
		}

		respRefund, _ := s.TopupWallet(refundWallet)
		txn.BalanceAfter = respRefund.Balance
	}

	err := s.txnRepository.UpdateTransaction(&txn)
	if err != nil {
		s.ctx.Logger().Warnf("failed to update transaction: %v", err)
		return nil
	}

	return nil
}
