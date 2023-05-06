package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	RefId           string
	UserId          string
	TransactionId   string
	CustomerNo      string
	BuyerSkuCode    string
	PricePaid       int64
	PriceDist       int64
	Status          string
	ResponseMessage string
	ResponseCode    string
	Sn              string
	BalanceBefore   int64
	BalanceAfter    int64
}
