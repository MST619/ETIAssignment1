package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var trips map[string]tripInfo

type tripInfo struct {
	Title string `json:"Trips"`
}

type Trip struct {
	TripID      int    `json:"tripid"`
	PickupPC    string `json:"pickuppc"`
	DropoffPC   string `json:"dropoffpc"`
	DriverID    int    `json:"driverid"`
	PassengerID int    `json:"passengerid"`
	TripStatus  string `json:"tripstatus"`
}

func tvalidKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func triphome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API for Trips!")
}

func alltrips(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, v)
	}
	json.NewEncoder(w).Encode(trips)
}

func trip(w http.ResponseWriter, r *http.Request) {
	if !tvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	//Database Connection
	if r.Header.Get("Content-type") == "application/json" {
		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Trips database opened!")
		}

		//THE GET REQUEST FOR TRIPS
		if r.Method == "GET" {
			params := mux.Vars(r)
			//fmt.Println(params)
			var getAllTrips Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &getAllTrips)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("Invalid trip!"))
					return
				}
			}
			json.NewEncoder(w).Encode(GetAllTripsRecord(db, params["tripid"]))
			w.WriteHeader(http.StatusAccepted)
			return
		}

		//POST request for creating new trip
		if r.Method == "POST" {
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else {
					if newTrip.PassengerID == 0 || newTrip.DriverID == 0 {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply trip " + "information " + "in JSON format"))
						return
					} else {
						if validateTrips(db, newTrip.PassengerID, newTrip.DriverID) {
							InsertTripRecord(db, newTrip.TripID, newTrip.PickupPC, newTrip.DropoffPC, newTrip.DriverID, newTrip.PassengerID, newTrip.TripStatus)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("201 - Trip added!"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("409 - Duplicate trip ID"))
							return
						}
					}
				}
			}
			//PUT for creating or updating existing trips
		} else if r.Method == "PUT" {
			//params := mux.Vars(r)
			var updateTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &updateTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else {
					if validateTrips(db, updateTrip.DriverID, updateTrip.PassengerID) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("No trip found!"))
					} else {
						if EditTripRecord(db, updateTrip.TripID, updateTrip.PickupPC, updateTrip.DropoffPC, updateTrip.DriverID, updateTrip.PassengerID, updateTrip.TripStatus) {
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("201 - Trip updated!"))
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("401 - Trip not updated!"))
						}
					}
				}
			}
		}
		//DELETE request but should not really work since you can't delete trips
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("404 - You are not able to delete your trip due to audit purposes"))
		}
	}
}

//Validate the trips
func validateTrips(db *sql.DB, PID int, DID int) bool {
	query := fmt.Sprintf("SELECT * FROM ETIAsgn.Trips WHERE PassengerID= '%d' OR DriverID='%d'", PID, DID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		err = results.Scan(&trip.TripID, &trip.PickupPC, &trip.DropoffPC, &trip.DriverID, &trip.PassengerID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		} else if trip.TripStatus != "Finished" {
			return false
		}
	}
	return true
}

//Get any free drivers who have 'Finished' their trip
func GetFreeDriver(db *sql.DB) int {
	query := fmt.Sprintf("SELECT Drivers.DriverID FROM Drivers INNER JOIN Trips ON Drivers.DriverID = Trips.DriverID WHERE Trips.TripStatus='Finished'")
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var DriverID int
	for results.Next() {
		err = results.Scan(&DriverID)
		if err != nil {
			panic(err.Error())
		}
	}
	return DriverID
}

//Get all the trips in the database
func GetAllTripsRecord(db *sql.DB, TID string) Trip {
	//query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID=?", TID)
	fmt.Println(TID)
	results, err := db.Query("SELECT * FROM Trips WHERE TripID=?", TID)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		err = results.Scan(&trip.TripID, &trip.PickupPC, &trip.DropoffPC, &trip.DriverID, &trip.PassengerID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		}
	}
	//fmt.Println(trip)
	return trip
}

//Get the specific trips the passenger has been on
func GetAllTrips(db *sql.DB, PID int) []Trip {
	query := fmt.Sprintf("SELECT * FROM Trips WHERE PassengerID ='%d'", PID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trips []Trip
	for results.Next() {
		var trip Trip
		err = results.Scan(&trip.TripID, &trip.PickupPC, &trip.DropoffPC, &trip.DriverID, &trip.PassengerID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		}
		trips = append(trips, trip)
	}
	return trips
}

//Insert a new trip into the database
func InsertTripRecord(db *sql.DB, TID int, PKPC string, DRPC string, DID int, PID int, TST string) bool {
	//query := fmt.Sprintf("INSERT INTO Trips VALUES (?,'%s','%s','%d','%d')", TID, PKPC, DRPC, DID, PID, TST)
	_, err := db.Query("INSERT INTO Trips VALUES (?, ?, ?, ?, ?, ?)", TID, PKPC, DRPC, DID, PID, TST)
	if err != nil {
		panic(err.Error())
	}
	return true
}

//Edit existing trip in the database
func EditTripRecord(db *sql.DB, TID int, PKPC string, DRPC string, DID int, PID int, TST string) bool {
	query := fmt.Sprintf("UPDATE Trips SET TripID='%d', PickupPC='%s', DropoffPC='%s', DriverID='%d', PassengerID='%d', TripStatus='%s' WHERE TripID='%d'", TID, PKPC, DRPC, DID, PID, TST)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

//Delete any data from the database. Should not really work since you can't delete trips
func DeleteTripRecord(db *sql.DB, TID int) {
	fmt.Println("Sorry. You are not able to delete your account due to audit purposes.")
}

//Main function for handling URL paths and handlers
func main() {
	trips = make(map[string]tripInfo)
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", triphome)
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods("GET", "PUT", "POST", "DELETE")
	//router.HandleFunc("/api/v1/trips", alltrips)

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))
}
