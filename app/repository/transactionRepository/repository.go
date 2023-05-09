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
		GetUserBalance(user_id string) dto.GetUserBalance
		UpdateWallet(wallet *model.Wallet) error
		GetWalletByUserID(wallet *model.Wallet) error
		GetProductByProductCode(product_code string) model.Product
		SaveTransaction(txn *model.Transaction) error
		UpdateTransaction(txn *model.Transaction) error
		GetTransactionByRefID(ref_id string) model.Transaction
		GetTransactionHistoryByUserID(user_id string) ([]model.Transaction, error)
		GetTransactionByTrxID(trx_id string) model.Transaction
	}
)

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) GetUserBalance(user_id string) dto.GetUserBalance {
	userBalance := dto.GetUserBalance{}

	r.db.Model(&model.User{}).Select("users.name, users.username, wallets.balance").Joins("join wallets on wallets.user_id = users.id").Where("users.uuid = ?", user_id).Limit(1).Find(&userBalance)

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
	txn := model.Transaction{}

	r.db.First(&txn, "ref_id = ?", ref_id)

	return txn
}

func (r *transactionRepository) GetTransactionByTrxID(trx_id string) model.Transaction {
	txn := model.Transaction{}

	r.db.First(&txn, "transaction_id = ?", trx_id)

	return txn
}

func (r *transactionRepository) GetTransactionHistoryByUserID(user_id string) ([]model.Transaction, error) {
	txn := []model.Transaction{}

	result := r.db.Where("user_id = ?", user_id).Find(&txn)

	if result.Error != nil {
		return txn, result.Error
	}

	return txn, nil
}
