package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API!")
}

func allPassengers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "List of all passengers")

	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
}

func passengers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Fprintf(w, "Passenger details"+params["passengerid"])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, r.Method)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)

	router.HandleFunc("/api/v1/passengers", allPassengers)
	router.HandleFunc("/api/v1/passengers/{passengerid}", passengers).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
