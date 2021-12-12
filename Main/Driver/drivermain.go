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

// used for storing drivers on the REST API
var drivers map[string]driverInfo

type driverInfo struct {
	Title string `json:"Driver"`
}

//Collections of fields for Driver and also to map this type to the record in the table
type Drivers struct {
	DriverID    string `json:"DriverID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber int    `json:"PhoneNumber"`
	Email       string `json:"Email"`
	LicenseNo   int    `json:"LicenseNo"`
}

//Access token used for securing the REST API
func dvalidKey(r *http.Request) bool {
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

func dhome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Driver REST API!")
}

//Main func for driver to call the requests, GET, PUT, POST, and DELETE
func driver(w http.ResponseWriter, r *http.Request) {

	if !dvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Driver database opened!")
	}

	//THE GET REQUEST FOR DRIVER
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getDrivers Drivers
		reqBody, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getDrivers)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetDriverRecords(db, params["driverid"], params["email"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {
		//POST for creating new driver
		if r.Method == "POST" {
			var newDriver Drivers
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newDriver)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("Error in JSON encoding. Error is %s", err)
				} else {
					if newDriver.Email == "" {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply driver information " + "in JSON format"))
						return
					} else {
						if !validateDriverRecord(db, newDriver.Email) {
							InsertDriverRecord(db, newDriver)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("201 - Driver added!"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("409 - Duplicate email"))
							return
						}
					}
				}
			}
			//PUT for creating or updating existing drivers
		} else if r.Method == "PUT" {
			var updateDriver Drivers
			reqbody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqbody, &updateDriver)

				if updateDriver.FirstName == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information " + "information " + "in JSON format"))
					return
				} else {
					if !validateDriverRecord(db, updateDriver.Email) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("No driver found with: " + updateDriver.Email))
					} else {
						EditDriverRecord(db, updateDriver)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Driver updated!"))
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

func GetAllDriverRecords(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(GetFreeDriver(db)))
}

func validateDriver(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["email"] == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply driver email " + "information " + "in JSON format"))
		return
	} else if validateDriverRecord(db, params["email"]) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(validateDriverEmail(db, params["email"])))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func GetDriverID(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn")
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["driverid"] == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply driver ID " + "information " + "in JSON format"))
		return
	} else {
		println(params["driverid"])
		query := fmt.Sprintf("SELECT LicenseNo FROM Drivers WHERE DriverID='%s'", params["driverid"])
		results, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		var LicenseNo string
		for results.Next() {
			err = results.Scan(&LicenseNo)
			if err != nil {
				panic(err.Error())
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(LicenseNo))
		return
	}
}

func validateDriverEmail(db *sql.DB, EML string) string {
	query := fmt.Sprintf("SELECT * FROM Drivers WHERE Email= '%s'", EML)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Drivers
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver.DriverID
}

func validateDriverRecord(db *sql.DB, EML string) bool {
	query := fmt.Sprintf("SELECT * FROM Drivers WHERE Email= '%s'", EML)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Drivers
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		} else if driver.Email == EML {
			return true
		}
	}
	return false
}

func GetDriverRecords(db *sql.DB, DID string, EML string) Drivers {
	//query := fmt.Sprintf("SELECT * FROM Drivers WHERE DriverID= ?", DID)
	results, err := db.Query("SELECT * FROM Drivers WHERE DriverID=? AND Email=?", DID, EML)
	if err != nil {
		panic(err.Error())
	}
	var driver Drivers
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver
}

func InsertDriverRecord(db *sql.DB, driver Drivers) bool {
	query := fmt.Sprintf("INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ('%s','%s','%s','%d','%s','%d');",
		driver.DriverID, driver.FirstName, driver.LastName, driver.PhoneNumber, driver.Email, driver.LicenseNo)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditDriverRecord(db *sql.DB, driver Drivers) bool {
	query := fmt.Sprintf("UPDATE Drivers SET FirstName='%s', LastName='%s', PhoneNumber=%d, Email='%s', LicenseNo='%d' WHERE DriverID='%s'",
		driver.FirstName, driver.LastName, driver.PhoneNumber, driver.Email, driver.LicenseNo, driver.DriverID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func GetFreeDriver(db *sql.DB) string {
	query := "SELECT DriverID FROM Drivers"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var DID string
	for results.Next() {
		var ID string
		err = results.Scan(&ID)
		if err != nil {
			panic(err.Error())
		}
		DID += ID + ","
	}
	return DID
}

func main() {
	drivers = make(map[string]driverInfo)
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/", dhome)
	router.HandleFunc("/api/v1/GetAllDriverRecords", GetAllDriverRecords)
	router.HandleFunc("/api/v1/validateDriverRecord/{email}", validateDriver)
	router.HandleFunc("/api/v1/GetDriver/{driverid}", GetDriverID)
	router.HandleFunc("/api/v1/drivers/{driverid}/{email}", driver).Methods("GET", "PUT", "POST", "DELETE")
	//router.HandleFunc("/api/v1/drivers", alldrivers)

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", handlers.CORS(headers, methods, origins)(router)))
}
