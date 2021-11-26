package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type tripInfo struct {
	Title string `json:"Trip"`
}

var trips map[string]tripInfo

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

func triphome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API for Trips!")
}

func alltrips(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "List of all trips")

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

	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)

	if r.Method == "GET" {
		if _, ok := trips[params["tripid"]]; ok {
			json.NewEncoder(w).Encode(trips[params["tripid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No trip found"))
		}
	}

	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("404 - You are not able to delete your account due to audit purposes"))
	}

	if r.Header.Get("Content-type") == "application/json" {
		//POST for creating new driver
		if r.Method == "POST" {
			var newTrip tripInfo
			reqbody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqbody, &newTrip)

				if newTrip.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply trip " + "information " + "in JSON format"))
					return
				}
				if _, ok := trips[params["tripid"]]; !ok {
					trips[params["tripid"]] = newTrip
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Trip added: " + params["tripid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate Trip ID"))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply trip information " + "in JSON format"))
			}
		}

		//PUT for creating or updating existing trips
		if r.Method == "PUT" {
			var newTrip tripInfo
			reqbody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqbody, &newTrip)

				if newTrip.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply trip information " + "information " + "in JSON format"))
					return
				}

				//check if passenger exists; add only if passenger does not exist
				if _, ok := trips[params["tripid"]]; !ok {
					trips[params["tripid"]] = newTrip
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Trip added: " + params["tripid"]))
				} else {
					//update trip
					trips[params["tripid"]] = newTrip
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Trip updated: " + params["trip"]))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply trip information " + "in JSON format"))
			}
		}
	}
}

func main() { //Even tho have error, it works....
	trips = make(map[string]tripInfo)
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", triphome)

	router.HandleFunc("/api/v1/trips", alltrips)
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
