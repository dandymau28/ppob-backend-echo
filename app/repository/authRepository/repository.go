package authRepository

import (
	"ppob-backend/app/dto"
	"ppob-backend/model"

	"gorm.io/gorm"
)

type (
	authRepository struct {
		db *gorm.DB
	}

	AuthRepository interface {
		GetUserCredentialByUsername(username string) dto.UserCredential
		SaveUser(user *model.User) error
		CreateWallet(user *model.User) error
	}
)

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) GetUserCredentialByUsername(username string) dto.UserCredential {
	user := dto.UserCredential{}

	r.db.Model(&model.User{}).Where("username = ?", username).Find(&user)
	return user
}

func (r *authRepository) SaveUser(user *model.User) error {
	result := r.db.Create(user)

	return result.Error
}

func (r *authRepository) CreateWallet(user *model.User) error {
	wallet := model.Wallet{
		UserID:  uint64(user.ID),
		Balance: 0,
	}

	result := r.db.Create(&wallet)

	return result.Error
}
