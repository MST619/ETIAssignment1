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

type Drivers struct {
	DriverID    int
	FirstName   string
	LastName    string
	PhoneNumber int
	Email       string
	LicenseNo   int
}

func getPassengerRecords(db *sql.DB) {
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

func getDriverRecords(db *sql.DB) {
	results, err := db.Query("Select * FROM ETIAsgn.Drivers")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var driver Drivers
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo)

		if err != nil {
			panic(err.Error())
		}
		fmt.Println(driver.DriverID, driver.FirstName, driver.LastName, driver.PhoneNumber, driver.Email, driver.LicenseNo)
	}
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Passengers")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}
	getPassengerRecords(db)
	getDriverRecords(db)
	defer db.Close()
}
