package model

import "gorm.io/gorm"

type User struct {
	Name     string
	Username string
	Password string
	Uuid     string
	gorm.Model
}
