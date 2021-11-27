package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type passengerInfo struct {
	Title string `json:"Passenger"`
}

func validKey(r *http.Request) bool {
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

var passengers map[string]passengerInfo

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API!")
}

func allPassengers(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "List of all passengers")

	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
	//returns all the passengers in JSON
	json.NewEncoder(w).Encode(passengers)
}

func passenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// fmt.Fprintf(w, "Passenger details"+params["passengerid"])
	// fmt.Fprintf(w, "\n")
	// fmt.Fprintf(w, r.Method)

	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	if r.Method == "GET" {
		if _, ok := passengers[params["passengerid"]]; ok {
			json.NewEncoder(w).Encode(
				passengers[params["passengerid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No passenger found"))
		}
	}

	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("404 - You are not able to delete your account due to audit purposes"))
	}

	if r.Header.Get("Content-type") == "application/json" {

		//POST for creating new passenger
		if r.Method == "POST" {
			var newPassenger passengerInfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
					return
				}

				if _, ok := passengers[params["passengerid"]]; !ok {
					passengers[params["passengerid"]] = newPassenger
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Passenger added: " + params["passengerid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate passenger ID"))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " + "in JSON format"))
			}
		}

		//PUT for creating or updating existing passengers
		if r.Method == "PUT" {
			var newPassenger passengerInfo
			reqbody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqbody, &newPassenger)

				if newPassenger.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger information " + "information " + "in JSON format"))
					return
				}
				// check if passenger exists; add only if passenger does not exist
				if _, ok := passengers[params["passengerid"]]; !ok {
					passengers[params["passengerid"]] = newPassenger
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Passenger added: " + params["passengerid"]))
				} else {
					//update passenger
					passengers[params["passengerid"]] = newPassenger
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Passenger updated: " + params["passengerid"]))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " + "information " + "in JSON format"))
			}
		}
	}
}

func main() {
	passengers = make(map[string]passengerInfo)
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)

	router.HandleFunc("/api/v1/passengers", allPassengers)
	router.HandleFunc("/api/v1/passengers/{passengerid}", passenger).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
