package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func Connect(connectionString string) error {
	var err error
	Db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Println("error in database.Connect")
		return err
	}
	log.Println("Connection DBName is Successful")
	return nil

}
