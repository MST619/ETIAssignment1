package main

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// type tripInfo struct {
// 	Title string `json:"Trip"`
// }

// var trips map[string]tripInfo

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API for Trips!")
}

// func alltrips(w http.ResponseWriter, r*http.Request){
// 	fmt.Fprintf(w, "List of all trips")

// 	kv := r.URL.Query()

// 	for k, v := range kv{
// 		fmt.Println(k,v)
// 	}
// 	json.NewEncoder(w).Encode(trips)
// }

func main() {
	//trips = make(map[string]tripInfo)
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)

	// router.HandleFunc("/api/v1/trips", alltrips)
	// router.HandleFunc("/api/v1/drivers/{tripid}", trip).Methods(
	// 	"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
