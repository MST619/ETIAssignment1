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

// used for storing passengers on the REST API
var passengers map[string]passengerInfo

func phome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Passenger REST API!")
}

type passengerInfo struct {
	Title string `json:"Passenger"`
}

//Collections of fields for Passengers and also to map this type to the record in the table
type Passengers struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber int    `json:"PhoneNumber"`
	Email       string `json:"Email"`
}

//Access token used for securing the REST API
func pvalidKey(r *http.Request) bool {
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

//Main func for passenger to call the requests, GET, PUT, POST, and DELETE
func passenger(w http.ResponseWriter, r *http.Request) {
	if !pvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}

	//THE GET request for passenger to retrive data from the Database.
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllPassengers Passengers
		reqBody, err := ioutil.ReadAll(r.Body)

		// defer the close till after the main function has finished executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllPassengers)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetPassengerRecord(db, params["passengerid"], params["email"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {
		//POST for creating new passenger
		if r.Method == "POST" {
			var newPassenger Passengers
			reqBody, err := ioutil.ReadAll(r.Body)

			// defer the close till after the main function has finished executing
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newPassenger)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else { //Check if passenger's email is empty or not
					if newPassenger.Email == "" {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
						return
					} else { //Validate passenger.
						if !validatePassengerRecord(db, newPassenger.Email) {
							InsertPassengerRecord(db, newPassenger)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("201 - Passenger added!"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("409 - Duplicate passenger ID"))
							return
						}
					}
				}
			}
			//PUT for creating or updating existing passengers
		} else if r.Method == "PUT" {
			fmt.Println("put called")
			var updatePassenger Passengers
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &updatePassenger)

				//Checking if the passenger's firstname is empty
				if updatePassenger.FirstName == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information " + "information " + "in JSON format"))
					return
				} else { //Checking to see if there is a existing passenger in the database
					if !validatePassengerRecord(db, updatePassenger.Email) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("No passenger found with: " + updatePassenger.Email))
					} else {
						EditPassengerRecord(db, updatePassenger.PassengerID, updatePassenger.FirstName, updatePassenger.LastName, updatePassenger.PhoneNumber, updatePassenger.Email)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Passenger updated!"))
						return
					}
				}
			}
		}
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("404 - You are not able to delete your account due to audit purposes"))
		}
	}
}

func allPassengers(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, v)
	}
	//returns all the passengers in JSON
	json.NewEncoder(w).Encode(passengers)
}

//To check if whether there is a duplicate email in the system
func validatePassengerRecord(db *sql.DB, EML string) bool {
	query := fmt.Sprintf("SELECT * FROM Passengers WHERE Email= '%s'", EML)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var passenger Passengers
	for results.Next() {
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.PhoneNumber, &passenger.Email)
		if err != nil {
			panic(err.Error())
		} else if passenger.Email == EML {
			return true
		}
	}
	return false
}

//Function to validate whether a specific Passenger exists.
func validatePassengerID(db *sql.DB, PID string) int {
	query := fmt.Sprintf("SELECT * FROM Passengers WHERE PassengerID=%s", PID)
	var passenger Passengers
	row := db.QueryRow(query) //Method to execute the query and is expected to return a single row.
	if err := row.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.PhoneNumber, &passenger.Email); err != nil {
		panic(err.Error())
	} else {
		return passenger.PassengerID
	}
}

func validatePassenger(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if _, err := strconv.Atoi(params["id"]); err != nil { //Converting string to int
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply passenger information " + "information " + "in JSON format"))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(validatePassengerID(db, params["id"])))) //Converting int to string
	}
}

func GetPassengerRecord(db *sql.DB, PID string, EML string) Passengers {
	results, err := db.Query("SELECT * FROM Passengers WHERE PassengerID=? AND Email=?", PID, EML)
	if err != nil {
		panic(err.Error())
	}
	var passenger Passengers
	for results.Next() {
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.PhoneNumber, &passenger.Email)
		if err != nil {
			panic(err.Error())
		}
	}
	return passenger
}

func InsertPassengerRecord(db *sql.DB, passenger Passengers) bool {
	query := fmt.Sprintf("INSERT INTO Passengers (PassengerID, FirstName, LastName, PhoneNumber, Email) VALUES ('%d','%s','%s','%d','%s');",
		passenger.PassengerID, passenger.FirstName, passenger.LastName, passenger.PhoneNumber, passenger.Email)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditPassengerRecord(db *sql.DB, PID int, FN string, LN string, PN int, EML string) bool {
	query := fmt.Sprintf("UPDATE Passengers SET FirstName='%s', LastName='%s', PhoneNumber=%d, Email='%s' WHERE PassengerID=%d", FN, LN, PN, EML, PID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func DeletePassengers(db *sql.DB, PID int) {
	fmt.Println("Sorry. You are not able to delete your account due to audit purposes.")
}

func main() {

	passengers = make(map[string]passengerInfo)
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/", phome)
	router.HandleFunc("/api/v1/validatePassengerRecord/{id}", validatePassenger)
	router.HandleFunc("/api/v1/passengers/{passengerid}/{email}", passenger).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(router)))
}
