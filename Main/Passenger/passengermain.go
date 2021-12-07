package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var passengers map[string]passengerInfo

func phome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Passenger REST API!")
}

type passengerInfo struct {
	Title string `json:"Passenger"`
}

//Collections of fields for Passengers
type Passengers struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber int    `json:"PhoneNumber"`
	Email       string `json:"Email"`
}

func pvalidKey(r *http.Request) bool {
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

	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllPassengers Passengers
		reqBody, err := ioutil.ReadAll(r.Body)
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
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newPassenger)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else {
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

				if updatePassenger.FirstName == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information " + "information " + "in JSON format"))
					return
				} else {
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

func GetPassengerRecord(db *sql.DB, PID string, EML string) Passengers {
	//query := fmt.Sprintf("SELECT * FROM ETIAsgn.Passengers WHERE PassengerID=?", PID)
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
	router.HandleFunc("/api/v1/passengers/{passengerid}/{email}", passenger).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(router)))
}
