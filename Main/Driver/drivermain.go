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

var drivers map[string]driverInfo

type driverInfo struct {
	Title string `json:"Driver"`
}

type Drivers struct {
	DriverID    int    `json:"DriverID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber int    `json:"PhoneNumber"`
	Email       string `json:"Email"`
	LicenseNo   int    `json:"LicenseNo"`
}

func dvalidKey(r *http.Request) bool {
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

// func alldrivers(w http.ResponseWriter, r *http.Request) {
// 	kv := r.URL.Query()
// 	for k, v := range kv {
// 		fmt.Println(k, v)
// 	}
// 	json.NewEncoder(w).Encode(drivers)
// }

//Main func for driver to call the requests, GET, PUT, POST, and DELETE
func driver(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	// fmt.Fprintf(w, "Detail for driver "+params["driverid"])
	// fmt.Fprintf(w, "\n")
	// fmt.Fprintf(w, r.Method)

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
		var getAllDrivers Drivers
		reqBody, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllDrivers)
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
						EditDriverRecord(db, updateDriver.DriverID, updateDriver.FirstName, updateDriver.LastName, updateDriver.PhoneNumber, updateDriver.Email, updateDriver.LicenseNo)
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

func validateDriverRecord(db *sql.DB, EML string) bool {
	query := fmt.Sprintf("SELECT * FROM ETIAsgn.Drivers WHERE Email= '%s'", EML)
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
	query := fmt.Sprintf("INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ('%d','%s','%s','%d','%s','%d');",
		driver.DriverID, driver.FirstName, driver.LastName, driver.PhoneNumber, driver.Email, driver.LicenseNo)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditDriverRecord(db *sql.DB, DID int, FN string, LN string, PN int, EML string, LCN int) bool {
	query := fmt.Sprintf("UPDATE Drivers SET FirstName='%s', LastName='%s', PhoneNumber=%d, Email='%s', LicenseNo='%d' WHERE DriverID='%d'",
		FN, LN, PN, EML, LCN, DID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func main() {
	drivers = make(map[string]driverInfo)
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/", dhome)
	router.HandleFunc("/api/v1/drivers/{driverid}/{email}", driver).Methods("GET", "PUT", "POST", "DELETE")
	//router.HandleFunc("/api/v1/drivers", alldrivers)

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", handlers.CORS(headers, methods, origins)(router)))
}
