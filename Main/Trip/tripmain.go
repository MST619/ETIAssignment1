package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// used for storing trips on the REST API
var trips map[string]tripInfo

type tripInfo struct {
	Title string `json:"Trips"`
}

//Collections of fields for Trip and also to map this type to the record in the table
type Trip struct {
	TripID      int    `json:"tripid"`
	PickupPC    string `json:"pickuppc"`
	DropoffPC   string `json:"dropoffpc"`
	DriverID    int    `json:"driverid"`
	PassengerID int    `json:"passengerid"`
	TripStatus  string `json:"tripstatus"`
}

//Collections of fields for Passengers and also to map this type to the record in the table
type Passengers struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber int    `json:"PhoneNumber"`
	Email       string `json:"Email"`
}

//Collections of fields for Drivers and also to map this type to the record in the table
type Drivers struct {
	DriverID    int    `json:"DriverID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber int    `json:"PhoneNumber"`
	Email       string `json:"Email"`
	LicenseNo   int    `json:"LicenseNo"`
}

//Access token used for securing the REST API
func tvalidKey(r *http.Request) bool {
	// returns the key/value pairs in the query string as a map object
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
	// Establishing the database Connection
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Trips database opened!")
	}

	//THE GET REQUEST FOR TRIPS
	if r.Method == "GET" {
		params := mux.Vars(r)
		fmt.Println(params)
		Tripid, err := strconv.Atoi(params["tripid"])
		if err != nil {
			fmt.Println(err)
		}
		tripInfo := GetAllTripsRecord(db, Tripid)
		if err != nil {
			fmt.Printf("Error in JSON encoding. Error is %s", err)
		} else if tripInfo.TripID == 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Invalid trip!"))
			return
		} else {
			json.NewEncoder(w).Encode(GetAllTripsRecord(db, tripInfo.TripID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		//POST request for creating new trip
		if r.Method == "POST" {
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			// defer the close till after the main function has finished executing
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
					}

					if validateTrips(db, newTrip.PassengerID, newTrip.DriverID) {
						InsertTripRecord(db, newTrip)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Trip added!"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("409 - Ongoing trip"))
						return
					}
				}
			}
			//PUT for creating or updating existing trips
		} else if r.Method == "PUT" {
			var updateTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			// defer the close till after the main function has finished executing
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &updateTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else {
					if validateTrips(db, updateTrip.DriverID, updateTrip.PassengerID) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("404 - No trip found!"))
					} else {
						if EditTripRecord(db, updateTrip) {
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

func validatePassengerRecord(PID int) int {
	url := "http://localhost:5000/api/v1/validatePassengerRecord/" + strconv.Itoa(PID)
	reposnse, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return 0
	}

	if err != nil {
		log.Fatal(err)
	} else if reposnse.StatusCode == http.StatusCreated {
		reposnseData, err := ioutil.ReadAll(reposnse.Body)
		if err != nil {
			println(err)
		} else {
			info, err := strconv.Atoi(string(reposnseData))
			if err != nil {
				println(err)
			}
			return info
		}
	}
	return 0
}

//Calling
func GetAllDriverRecords() Drivers {
	response, err := http.Get("http://localhost:5001/api/v1/GetAllDriverRecords/")
	if err != nil {
		fmt.Print(err.Error())
	}
	var driverTrip Drivers
	if response.StatusCode == http.StatusAccepted {
		responseData, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(responseData))
		fmt.Println(response.StatusCode)
		response.Body.Close()
		json.Unmarshal(responseData, &driverTrip)
	} else {
		fmt.Printf("404 - There are no free drivers at the moment")
		return driverTrip
	}
	return driverTrip
}

func validateDriver(Email string) string {
	url := "http://localhost:5001/api/v1/validateDriverRecord/" + Email
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			return string(responseData)
		}
	}
	return ""
}

func GetDriver(DriverID string) string {
	url := "http://localhost:5001/api/v1/GetDriver/" + DriverID

	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			return string(responseData)
		}
	}
	return ""
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
func GetFreeDriver(db *sql.DB, DID int) int {
	//query := fmt.Sprintf("SELECT Drivers.DriverID FROM Drivers INNER JOIN Trips ON Drivers.DriverID = Trips.DriverID WHERE Trips.TripStatus='Finished'")
	results, err := db.Query("SELECT Drivers.DriverID FROM Drivers INNER JOIN Trips ON Drivers.DriverID = Trips.DriverID WHERE Trips.TripStatus='Finished'")
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
func GetAllTripsRecord(db *sql.DB, TID int) Trip {
	//query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID=?", TID)
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
func InsertTripRecord(db *sql.DB, trip Trip) bool {
	query := fmt.Sprintf("INSERT INTO Trips (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ('%d','%s','%s','%d','%d','%s')",
		trip.TripID, trip.PickupPC, trip.DropoffPC, trip.DriverID, trip.PassengerID, trip.TripStatus)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

//Edit existing trip in the database
func EditTripRecord(db *sql.DB, trip Trip) bool {
	query := fmt.Sprintf("UPDATE Trips SET PickupPC='%s', DropoffPC='%s', DriverID='%d', PassengerID='%d', TripStatus='%s' WHERE TripID='%d'",
		trip.PickupPC, trip.DropoffPC, trip.DriverID, trip.PassengerID, trip.TripStatus, trip.TripID)
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
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/", triphome)
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods("GET", "PUT", "POST", "DELETE")
	router.HandleFunc("/api/v1/trips", alltrips)

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", handlers.CORS(headers, methods, origins)(router)))
}
