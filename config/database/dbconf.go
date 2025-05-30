package database

import (
	"log"
	"os"

	"github.com/LeandroDeJesus-S/confectionery/api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseStarter is an structure to manage the database startup processes
type DatabaseStarter struct {
	db *gorm.DB
}

// NewDatabaseStarter initializes the DatabaseStarter struct
func NewDatabaseStarter() *DatabaseStarter {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_STRING")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &DatabaseStarter{db: db}
}

// MakeMigrations performs all the migrations process
func (d *DatabaseStarter) MakeMigrations() {
	d.db.AutoMigrate(
		&models.Customer{},
		&models.Cake{},
		&models.Order{},
	)
}
