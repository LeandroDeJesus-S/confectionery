package models

// Cake stores data about a cake
type Cake struct {
	ID    uint  
	Name  string `gorm:"unique;not null;size:100"`
	Price uint64  `gorm:"not null;default:0"`
}
