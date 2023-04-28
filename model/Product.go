package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName         string
	Category            string
	Brand               string
	Type                string
	SellerName          string
	Price               int64
	SellerPrice         int64
	BuyerSkuCode        string
	BuyerProductStatus  bool
	SellerProductStatus bool
	UnlimitedStock      bool
	Stock               int
	Multi               bool
	StartCutOff         string
	EndCutOff           string
	Desc                string
	Uuid                string
}
