package orm

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name   string
	Status int
	ShopId int
}
