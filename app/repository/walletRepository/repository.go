package walletRepository

import (
	"ppob-backend/app/dto"
	"ppob-backend/model"

	"gorm.io/gorm"
)

type (
	walletRepository struct {
		db *gorm.DB
	}

	WalletRepository interface {
		GetUserBalance(userID string) dto.GetUserBalance
		UpdateWallet(wallet *model.Wallet) error
		GetWalletByUserID(wallet *model.Wallet) error
	}
)

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) GetUserBalance(user_id string) dto.GetUserBalance {
	userBalance := dto.GetUserBalance{}

	r.db.Model(&model.User{}).Select("users.name, users.username, wallets.balance").Joins("join wallets on wallets.user_id = users.id").Where("users.uuid = ?", user_id).Limit(1).Find(&userBalance)

	return userBalance
}

func (r *walletRepository) GetWalletByUserID(wallet *model.Wallet) error {
	result := r.db.First(wallet)

	return result.Error
}

func (r *walletRepository) UpdateWallet(wallet *model.Wallet) error {
	result := r.db.Save(wallet)

	return result.Error
}

func (r *walletRepository) GetProductByProductCode(product_code string) model.Product {
	product := model.Product{}

	r.db.Model(&model.Product{}).Where("buyer_sku_code", product_code).Find(&product)
	return product
}

func (r *walletRepository) SaveTransaction(txn *model.Transaction) error {
	result := r.db.Create(txn)

	return result.Error
}

func (r *walletRepository) UpdateTransaction(txn *model.Transaction) error {
	result := r.db.Save(txn)

	return result.Error
}

func (r *walletRepository) GetTransactionByRefID(ref_id string) model.Transaction {
	txn := model.Transaction{
		RefId: ref_id,
	}

	r.db.First(&txn)

	return txn
}

func (r *walletRepository) GetTransactionHistoryByUserID(user_id uint64) ([]model.Transaction, error) {
	txn := []model.Transaction{}

	result := r.db.Where("user_id = ?", user_id).Find(&txn)

	if result.Error != nil {
		return txn, result.Error
	}

	return txn, nil
}
