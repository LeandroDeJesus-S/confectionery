package models

import (
	"gorm.io/gorm"
)

// Order represents the schema of an order made by a customer
type Order struct {
	gorm.Model
	CustomerID uint  `gorm:"primaryKey"`
	CakeID     uint  `gorm:"primaryKey"`
	Qtd        uint  
	Delivered  bool  `gorm:"default:false"`
}
