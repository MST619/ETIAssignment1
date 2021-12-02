package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type driverInfo struct {
	Title string `json:"Driver"`
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

var drivers map[string]driverInfo

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API!")
}

func alldrivers(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "List of all drivers")

	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
	json.NewEncoder(w).Encode(drivers)
}

func driver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// fmt.Fprintf(w, "Detail for driver "+params["driverid"])
	// fmt.Fprintf(w, "\n")
	// fmt.Fprintf(w, r.Method)

	if !dvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	if r.Method == "GET" {
		if _, ok := drivers[params["driverid"]]; ok {
			json.NewEncoder(w).Encode(drivers[params["driverid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No driver found"))
		}
	}

	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("404 - You are not able to delete your account due to audit purposes"))
	}

	if r.Header.Get("Content-type") == "application/json" {
		//POST for creating new driver
		if r.Method == "POST" {
			var newDriver driverInfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver " + "information " + "in JSON format"))
					return
				}
				if _, ok := drivers[params["driverid"]]; !ok {
					drivers[params["driverid"]] = newDriver
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver added: " + params["driverid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate Driver ID"))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver information " + "in JSON format"))
			}
		}

		//PUT for creating or updating existing drivers
		if r.Method == "PUT" {
			var newDriver driverInfo
			reqbody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqbody, &newDriver)

				if newDriver.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information " + "information " + "in JSON format"))
					return
				}
				//check if passenger exists; add only if passenger does not exist
				if _, ok := drivers[params["driverid"]]; !ok {
					drivers[params["driverid"]] = newDriver
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver added: " + params["driverid"]))
				} else {
					//update passenger
					drivers[params["driverid"]] = newDriver
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Driver updated: " + params["driverid"]))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver information " + "in JSON format"))
			}
		}
	}
}

func main() {
	drivers = make(map[string]driverInfo)
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)

	router.HandleFunc("/api/v1/drivers", alldrivers)
	router.HandleFunc("/api/v1/drivers/{driverid}", driver).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
