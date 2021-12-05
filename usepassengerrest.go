package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// const baseURL = "http://localhost:5000/api/v1/passengers"
// const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

// func getPassenger(code string) {
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

// func addPassenger(code string, jsonData map[string]string) {
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

// func updatePassenger(code string, jsonData map[string]string) {
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

// func deletePassenger(code string) {
// 	fmt.Println("404 - You are not able to delete your passenger account due to audit purposes???")
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
// 	jsonData := map[string]string{"passenger": "Joke Peralta"}
// 	addPassenger("9544", jsonData)
// 	// 201
// 	// 201 - Passenger added: 9544

// 	jsonData = map[string]string{"passenger": "Amy Santiago"}
// 	addPassenger("3263", jsonData)
// 	// 201
// 	// 201 - Passenger added 3263

// 	jsonData = map[string]string{"passenger": "Jake Peralta"}
// 	updatePassenger("9544", jsonData)
// 	// 202
// 	// 202 – Passengers updated: 9544

// 	getPassenger("") // get all Passengers
// 	// 200
// 	//
// 	//

// 	getPassenger("9544") // get a specific passenger
// 	// 200
// 	// {"Title":"Jake Peralta"}

// 	//deletePassenger("3263")
// 	//202
// 	//202 – Passenger deleted: 3263

// 	getPassenger("") // get all passengers
// 	// 200
// 	// {Both passengers should still be here}
// }
