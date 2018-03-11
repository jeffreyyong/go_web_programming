package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Every model(table) created should be represented as struct in GORM. And the first line should be gorm.Model
// The other fields are the fields of the table
// By default, an incrementing ID will be created
type User struct {
	gorm.Model
	Orders []Order
	Data   string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"-"`
}

type Order struct {
	gorm.Model
	User User
	Data string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}

// GORM creates tables with plural names. Use this to suppress it
func (User) TableName() string {
	return "user"
}

func (Order) TableName() string {
	return "order"
}

func InitDB() (*gorm.DB, error) {
	var err error
	db, err := gorm.Open("postgres", "postgres://jeffrey:passme123@localhost/mydb?sslmode=disable")
	if err != nil {
		return nil, err
	} else {
		// This function creates the tables for structs passed as parameters. It makes sure that if tables exist already, it skips creation
		db.AutoMigrate(&User{}, &Order{})
		return db, nil
	}
}
