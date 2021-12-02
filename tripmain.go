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
	PostalCode  int    `json:"PostalCode"`
	Pickup      string `json:"Pickup"`
	Dropoff     string `json:"Dropoff"`
	DriverID    int    `json:"DriverID"`
	PassengerID int    `json:"PassengerID"`
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
	// fmt.Fprintf(w, "Detail for trips "+params["tripid"])
	// fmt.Fprintf(w, "\n")
	// fmt.Fprintf(w, r.Method)

	if !tvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	//params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {
		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Trip database opened!")
		}

		//THE GET REQUEST FOR TRIPS
		if r.Method == "GET" {
			var getAllTrips Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &getAllTrips)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else if getAllTrips.Pickup != "" || getAllTrips.Dropoff != "" {
					json.NewEncoder(w).Encode(GetAllTripsRecord(db, getAllTrips.Pickup, getAllTrips.Dropoff))
					w.WriteHeader(http.StatusAccepted)
					return
				} else {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("Invalid information!"))
					return
				}
			}
		}

		//POST for creating new driver
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
					if newTrip.Pickup == "" {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply trip " + "information " + "in JSON format"))
						return
					} else {
						if !validateTrips(db, newTrip.PostalCode) {
							InsertTripRecord(db, newTrip.PostalCode, newTrip.Pickup, newTrip.Dropoff, newTrip.DriverID, newTrip.PassengerID)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("201 - Trip added!"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("409 - Duplicate postal code"))
							return
						}
					}
				}
			}
			//PUT for creating or updating existing trips
		} else if r.Method == "PUT" {
			var updateTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &updateTrip)

				if updateTrip.Pickup == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply trip information " + "information " + "in JSON format"))
					return
				} else {
					if !validateTrips(db, updateTrip.PostalCode) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("No trip found with: " + updateTrip.Pickup))
					} else {
						EditTripRecord(db, updateTrip.PostalCode, updateTrip.Pickup, updateTrip.Dropoff, updateTrip.DriverID, updateTrip.PassengerID)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Passenger updated!"))
					}
				}
			}
		}
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("404 - You are not able to delete your trip due to audit purposes"))
		}
	}
}

func validateTrips(db *sql.DB, PSC int) bool {
	query := fmt.Sprintf("SELECT * FROM ETIAsgn.Trip WHERE PostalCode= '%d'", PSC)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		err = results.Scan(&trip.PostalCode, &trip.Pickup, &trip.Dropoff, &trip.DriverID, &trip.PassengerID)
		if err != nil {
			panic(err.Error())
		} else if trip.PostalCode == PSC {
			return true
		}
	}
	return false
}

func GetAllTripsRecord(db *sql.DB, PKUP string, DRP string) Trip {
	query := fmt.Sprintf("SELECT * FROM ETIAsgn.Trip WHERE Pickup= '%s' AND Dropoff='%s'", PKUP, DRP)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		err = results.Scan(&trip.PostalCode, &trip.Pickup, &trip.Dropoff, &trip.DriverID, &trip.PassengerID)
		if err != nil {
			panic(err.Error())
		}
	}
	return trip
}

func InsertTripRecord(db *sql.DB, PSC int, PKUP string, DRP string, DID int, PID int) bool {
	query := fmt.Sprintf("INSERT INTO Trip VALUES ('%d','%s','%s','%d','%d')",
		PSC, PKUP, DRP, DID, PID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditTripRecord(db *sql.DB, PSC int, PKUP string, DRP string, DID int, PID int) bool {
	query := fmt.Sprintf("UPDATE Trip SET PostalCode='%d', Pickup='%s', Dropoff='%s', DriverID='%d', PassengerID='%d' WHERE PostalCode='%d'",
		PSC, PKUP, DRP, DID, PID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func DeleteTripRecord(db *sql.DB, PSC int) {
	fmt.Println("Sorry. You are not able to delete your account due to audit purposes.")
}

func main() {
	trips = make(map[string]tripInfo)
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", triphome)
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods("GET", "PUT", "POST", "DELETE")
	//router.HandleFunc("/api/v1/trips", alltrips)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
