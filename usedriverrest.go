package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// const baseURL = "http://localhost:5000/api/v1/drivers"
// const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

// func getDriver(code string) {
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

// func addDriver(code string, jsonData map[string]string) {
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

// func updateDriver(code string, jsonData map[string]string) {
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

// func deleteDriver(code string) {
// 	fmt.Println("404 - You are not able to delete your driver account due to audit purposes???")
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
// 	jsonData := map[string]string{"driver": "Min Verstappen"}
// 	addDriver("33", jsonData)
// 	// 201
// 	// 201 - Driver added: NS16

// 	jsonData = map[string]string{"driver": "Lewis Hamilton"}
// 	addDriver("44", jsonData)
// 	// 201
// 	// 201 - Driver added DT5

// 	jsonData = map[string]string{"driver": "Max Verstappen"}
// 	updateDriver("33", jsonData)
// 	// 202
// 	// 202 – Driver updated: 33

// 	getDriver("") // get all Drivers
// 	// 200
// 	//
// 	//

// 	getDriver("33") // get a specific Driver
// 	// 200
// 	// {"Title":"Max Verstappen"}

// 	//deleteDriver("44")
// 	//202
// 	//202 – Driver deleted: 44

// 	getDriver("") // get all drivers
// 	// 200
// 	// {Both drivers should still be here}
// }
