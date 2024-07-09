package model

type User struct {
	Name     string
	Username string
	Password string
	Uuid     string `gorm:"primarykey"`
	BaseModel
}
