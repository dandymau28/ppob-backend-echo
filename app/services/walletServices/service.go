package walletServices

import (
	"errors"
	"ppob-backend/app/dto"
	"ppob-backend/app/repository/walletRepository"
	"ppob-backend/config"
	"ppob-backend/model"
	"sync"
)

type (
	walletServices struct {
		walletRepository walletRepository.WalletRepository
		Config           *config.SystemConfig
		Mutex            sync.Mutex
	}

	WalletServices interface {
		GetUserBalanceByUserID(request dto.GetBalanceRequest) (dto.GetBalanceResponse, error)
		TopupWallet(request dto.TopupWalletRequest) (dto.TopupWalletResponse, error)
	}
)

func NewWalletServices(config *config.SystemConfig, walletRepo walletRepository.WalletRepository) WalletServices {
	return &walletServices{
		walletRepository: walletRepo,
		Config:           config,
	}
}

func (s *walletServices) GetUserBalanceByUserID(request dto.GetBalanceRequest) (dto.GetBalanceResponse, error) {
	var (
		response = dto.GetBalanceResponse{}
	)

	userID := request.UserID

	userBalance := s.walletRepository.GetUserBalance(userID)

	if (dto.GetUserBalance{}) == userBalance {
		return response, errors.New("no user balance found")
	}

	response = dto.GetBalanceResponse(userBalance)

	return response, nil
}

func (s *walletServices) TopupWallet(request dto.TopupWalletRequest) (dto.TopupWalletResponse, error) {
	var (
		response = dto.TopupWalletResponse{}
		wallet   = model.Wallet{}
	)

	userID := request.UserID

	wallet.Uuid = userID

	s.Mutex.Lock()
	err := s.walletRepository.GetWalletByUserID(&wallet)
	balanceBefore := wallet.Balance

	if err != nil {
		return response, errors.New("failed to get wallet")
	}

	wallet.Balance += int64(request.Amount)

	if wallet.Balance < 0 {
		s.Mutex.Unlock()
		return response, errors.New("insufficient balance")
	}

	err = s.walletRepository.UpdateWallet(&wallet)

	if err != nil {
		return response, errors.New("failed to update wallet")
	}
	s.Mutex.Unlock()

	userBalance := s.walletRepository.GetUserBalance(userID)

	response = dto.TopupWalletResponse{
		Name:          userBalance.Name,
		Username:      userBalance.Username,
		Balance:       userBalance.Balance,
		BalanceBefore: balanceBefore,
	}

	return response, nil
}
