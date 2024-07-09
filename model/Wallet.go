package model

type Wallet struct {
	BaseModel
	UserID  string
	Balance int64
	Uuid    string `gorm:"primarykey"`
}
