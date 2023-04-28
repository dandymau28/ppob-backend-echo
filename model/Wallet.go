package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID  uint64
	Balance int64
	Uuid    string
}
