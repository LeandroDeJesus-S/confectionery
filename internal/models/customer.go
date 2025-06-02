package models

type Customer struct {
	ID       uint
	Fname    string
	Lname    string
	Email    string `gorm:"unique"`
	Active   bool   `gorm:"default:true"`
	Orders []Order  `gorm:"many2many:orders;"`
}
