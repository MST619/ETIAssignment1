package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//Collections of fields for Passengers
type Passengers struct {
	PassengerID int
	FirstName   string
	LastName    string
	PhoneNumber int
	Email       string
}

//Collections of fields for Drivers
type Drivers struct {
	DriverID    int
	FirstName   string
	LastName    string
	PhoneNumber int
	Email       string
	LicenseNo   int
}

//Collections of fields for trips
type Trip struct {
	PostalCode  int
	Pickup      string
	Dropoff     string
	DriverID    int
	PassengerID int
}

//Function to get the passenger table from the database
func getPassengerRecords(db *sql.DB) {
	results, err := db.Query("SELECT * FROM ETIAsgn.Passengers")

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

//function to get the driver table from the database
// func getDriverRecords(db *sql.DB) {
// 	results, err := db.Query("SELECT * FROM ETIAsgn.Drivers")

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for results.Next() {
// 		var driver Drivers
// 		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo)

// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Println(driver.DriverID, driver.FirstName, driver.LastName, driver.PhoneNumber, driver.Email, driver.LicenseNo)
// 	}
// }

//function to get the trip table from the database
// func getTripRecords(db *sql.DB) {
// 	results, err := db.Query("SELECT * FROM ETIAsgn.Trip")

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for results.Next() {
// 		var trip Trip
// 		err = results.Scan(&trip.PostalCode, &trip.Pickup, &trip.Dropoff, &trip.DriverID, &trip.PassengerID)

// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Println(trip.PostalCode, trip.Pickup, trip.Dropoff, trip.DriverID, trip.PassengerID)
// 	}
// }

// func InsertPassengerRecord(db *sql.DB, PID int, FN string, LN string, PN int, EML string) {
// 	query := fmt.Sprintf("INSERT INTO Passengers VALUES (%d, '%s', '%s', %d, '%s')", PID, FN, LN, PN, EML)

// 	_, err := db.Query(query)

// 	if err != nil {
// 		panic(err.Error())
// 	}
// }

// func EditPassengeRecord(db *sql.DB, PID int, FN string, LN string, PN int, EML string) {
// 	query := fmt.Sprintf(
// 		"UPDATE Passengers SET FirstName='%s', LastName='%s', PhoneNumber=%d, Email='%s' WHERE PassengerID=%d",
// 		FN, LN, PN, EML, PID)
// 	_, err := db.Query(query)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }

func DeletePassengers(db *sql.DB, PID int) {
	fmt.Println("Sorry. You are not able to delete your account due to audit purposes.")
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}

	//InsertPassengerRecord(db, 6, "Rosa", "Diaz", 44443333, "rdb99@gmail.com")
	//EditPassengeRecord(db, 0006, "Rosa", "Diaz", 43343443, "rdb99@gmail.com")
	DeletePassengers(db, 6)

	getPassengerRecords(db)
	//getDriverRecords(db)
	//getTripRecords(db)
	defer db.Close()
}
