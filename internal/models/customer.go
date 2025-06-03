package models

type Customer struct {
	ID       uint
	Fname    string `gorm:"size:100;not null"`
	Lname    string `gorm:"size:255;not null"`
	Email    string `gorm:"unique;not null;size:345"`
	Active   bool   `gorm:"default:true"`
	Orders []Order  `gorm:"many2many:orders;"`
}
