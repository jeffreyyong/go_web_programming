package main

import (
	"log"

	"web/postgres/basicexample/models"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Printf("db %v error: %v", db, err)
	}
}
