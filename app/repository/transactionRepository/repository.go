package transactionRepository

import (
	"ppob-backend/app/dto"
	"ppob-backend/model"

	"gorm.io/gorm"
)

type (
	transactionRepository struct {
		db *gorm.DB
	}

	TransactionRepository interface {
		GetUserBalance(userID uint64) dto.GetUserBalance
		UpdateWallet(wallet *model.Wallet) error
		GetWalletByUserID(wallet *model.Wallet) error
		GetProductByProductCode(product_code string) model.Product
		SaveTransaction(txn *model.Transaction) error
		UpdateTransaction(txn *model.Transaction) error
		GetTransactionByRefID(ref_id string) model.Transaction
		GetTransactionHistoryByUserID(user_id uint64) ([]model.Transaction, error)
	}
)

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) GetUserBalance(user_id uint64) dto.GetUserBalance {
	userBalance := dto.GetUserBalance{}

	r.db.Model(&model.User{}).Select("users.name, users.username, wallets.balance").Joins("join wallets on wallets.user_id = users.id").Where("users.id = ?", user_id).Limit(1).Find(&userBalance)

	return userBalance
}

func (r *transactionRepository) GetWalletByUserID(wallet *model.Wallet) error {
	result := r.db.First(wallet)

	return result.Error
}

func (r *transactionRepository) UpdateWallet(wallet *model.Wallet) error {
	result := r.db.Save(wallet)

	return result.Error
}

func (r *transactionRepository) GetProductByProductCode(product_code string) model.Product {
	product := model.Product{}

	r.db.Model(&model.Product{}).Where("buyer_sku_code", product_code).Find(&product)
	return product
}

func (r *transactionRepository) SaveTransaction(txn *model.Transaction) error {
	result := r.db.Create(txn)

	return result.Error
}

func (r *transactionRepository) UpdateTransaction(txn *model.Transaction) error {
	result := r.db.Save(txn)

	return result.Error
}

func (r *transactionRepository) GetTransactionByRefID(ref_id string) model.Transaction {
	txn := model.Transaction{
		RefId: ref_id,
	}

	r.db.First(&txn)

	return txn
}

func (r *transactionRepository) GetTransactionHistoryByUserID(user_id uint64) ([]model.Transaction, error) {
	txn := []model.Transaction{}

	result := r.db.Where("user_id = ?", user_id).Find(&txn)

	if result.Error != nil {
		return txn, result.Error
	}

	return txn, nil
}
