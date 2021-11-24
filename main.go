package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Passengers struct {
	PassengerID int
	FirstName   string
	LastName    string
	PhoneNumber int
	Email       string
}

func getRecords(db *sql.DB) {
	results, err := db.Query("Select * FROM ETIAsgn.Passengers")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var passenger Passengers
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.PhoneNumber, &passenger.Email)

		if err != nil {
			panic(err.Error())
		}
		fmt.Println(passenger.PassengerID, passenger.FirstName, passenger.LastName, passenger.PhoneNumber, passenger.Email)
	}
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Passengers")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}
	getRecords(db)
	defer db.Close()
}
