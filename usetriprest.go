package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// const baseURL = "http://localhost:5000/api/v1/trips"
// const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

// func getTrip(code string) {
// 	url := baseURL
// 	if code != "" {
// 		url = baseURL + "/" + code + "?key=" + key
// 	}
// 	response, err := http.Get(url)

// 	if err != nil {
// 		fmt.Printf("The HTTP request failed with error %s\n", err)
// 	} else {
// 		data, _ := ioutil.ReadAll(response.Body)
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(string(data))
// 		response.Body.Close()
// 	}
// }

// func addTrip(code string, jsonData map[string]string) {
// 	jsonValue, _ := json.Marshal(jsonData)

// 	response, err := http.Post(baseURL+"/"+code+"?key="+key,
// 		"application/json", bytes.NewBuffer(jsonValue))

// 	if err != nil {
// 		fmt.Printf("The HTTP request failed with error %s\n", err)
// 	} else {
// 		data, _ := ioutil.ReadAll(response.Body)
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(string(data))
// 		response.Body.Close()
// 	}
// }

// func updateTrip(code string, jsonData map[string]string) {
// 	jsonValue, _ := json.Marshal(jsonData)

// 	request, err := http.NewRequest(http.MethodPut,
// 		baseURL+"/"+code+"?key="+key,
// 		bytes.NewBuffer(jsonValue))

// 	request.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	response, err := client.Do(request)

// 	if err != nil {
// 		fmt.Printf("The HTTP request failed with error %s\n", err)
// 	} else {
// 		data, _ := ioutil.ReadAll(response.Body)
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(string(data))
// 		response.Body.Close()
// 	}
// }

// func deleteTrip(code string) {
// 	fmt.Println("404 - You are not able to delete your trip due to audit purposes???")
// 	// request, err := http.NewRequest(http.MethodDelete,
// 	//     baseURL+"/"+code+"?key="+key, nil)

// 	// client := &http.Client{}
// 	// response, err := client.Do(request)

// 	// if err != nil {
// 	//     fmt.Printf("The HTTP request failed with error %s\n", err)
// 	// } else {
// 	//     data, _ := ioutil.ReadAll(response.Body)
// 	//     fmt.Println(response.StatusCode)
// 	//     fmt.Println(string(data))
// 	//     response.Body.Close()
// 	// }
// }

// func main() {
// 	jsonData := map[string]string{"trip": "Ang Mo Kio"}
// 	addTrip("NS16", jsonData)
// 	// 201
// 	// 201 - Trip added: NS16

// 	jsonData = map[string]string{"trip": "Bukit Timah"}
// 	addTrip("DT5", jsonData)
// 	// 201
// 	// 201 - Trip added DT5

// 	jsonData = map[string]string{"trip": "Angmo Kio"}
// 	updateTrip("NS16", jsonData)
// 	// 202
// 	// 202 – Trip updated: DT6

// 	getTrip("") // get all trips
// 	// 200
// 	//
// 	//

// 	getTrip("NS16") // get a specific trip
// 	// 200
// 	// {"Title":"Sir Lewis Hamilton"}

// 	//deleteCourse("IOT201")
// 	//202
// 	//202 – Course deleted: IOT201

// 	getTrip("") // get all courses
// 	// 200
// 	// {"44":{"Title":"Sir Lewis Hamilton"}}
// }
