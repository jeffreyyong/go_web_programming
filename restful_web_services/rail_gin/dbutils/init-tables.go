package dbutils

import (
	"database/sql"
	"log"
)

func Initialize(dbDriver *sql.DB) {
	statement, driverError := dbDriver.Prepare(train)
	if driverError != nil {
		log.Println(driverError)
	}
	// Create train table
	_, statementError := statement.Exec()
	if statementError != nil {
		log.Println("Table already exists!")
	}

	statement, _ = dbDriver.Prepare(station)
	statement.Exec()

	statement, err := dbDriver.Prepare(schedule)
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
	}
	log.Println("All tables created/initialized successfully!")
}
