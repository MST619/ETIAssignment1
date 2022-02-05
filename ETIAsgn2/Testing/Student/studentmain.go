package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// used for storing students on the REST API
var students map[string]studentInfo

func phome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the student REST API!")
}

type studentInfo struct {
	Title string `json:"Student"`
}

//Collections of fields for Passengers and also to map this type to the record in the table
type Students struct {
	StudentID   int    `json:"StudentID"`
	StudentName string `json:"StudentName"`
	DOB         string `json:"DOB"`
	Address     string `json:"Address"`
	PhoneNumber int    `json:"PhoneNumber"`
}

type Modules struct {
	ModuleCode      string `json:"ModuleCode"`
	ModuleName      string `json:"ModuleName"`
	ModuleSynopsis  string `json:"ModuleSynopsis"`
	ModuleObjective string `json:"ModuleObjective"`
	StudentID       int    `json:"StudentID"`
}

type Results struct {
	ResultsID    int    `json:"ResultsID"`
	ResultsGrade string `json:"ResultsGrade"`
	StudentID    int    `json:"StudentID"`
	ModuleCode   string `json:"ModuleCode"`
}

type Timetable struct {
	TimetableID string `json:"TimetableID"`
	LessonDay   string `json:"LessonDay"`
	StartTime   string `json:"StartTime"`
	EndTime     string `json:"EndTime"`
	ModuleCode  string `json:"ModuleCode"`
}

type CommentsRatings struct {
	RatingsID string `json:"RatingsID"`
	Ratings   string `json:"Ratings"`
	Comments  string `json:"Comments"`
	StudentID int    `json:"StudentID"`
}

//Access token used for securing the REST API
func pvalidKey(r *http.Request) bool {
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

//Main func for student to call the requests, GET, PUT, POST, and DELETE
func student(w http.ResponseWriter, r *http.Request) {
	if !pvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn2")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}

	//THE GET request for student to retrive data from the Database.
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllStudents Students
		reqBody, err := ioutil.ReadAll(r.Body)

		// defer the close till after the main function has finished executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllStudents)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetStudentRecord(db, params["studentid"], params["dob"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {
		//PUT for creating or updating existing students
		if r.Method == "PUT" {
			fmt.Println("put called")
			var updateStudent Students
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &updateStudent)

				//Checking if the student's name is empty
				if updateStudent.StudentName == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply student information " + "information " + "in JSON format"))
					return
				} else { //Checking to see if there is a existing student in the database
					if !validateStudentRecord(db, updateStudent.DOB) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("No student found with: " + updateStudent.DOB))
					} else {
						EditStudentRecord(db, updateStudent.StudentID, updateStudent.StudentName, updateStudent.DOB, updateStudent.Address, updateStudent.PhoneNumber)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Student updated!"))
						return
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

func module(w http.ResponseWriter, r *http.Request) {
	if !pvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn2")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}

	//THE GET request for student to retrive data from the Database.
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllModules Modules
		reqBody, err := ioutil.ReadAll(r.Body)

		// defer the close till after the main function has finished executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllModules)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetModules(db, params["modulecode"], params["studentid"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func results(w http.ResponseWriter, r *http.Request) {
	if !pvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn2")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}

	//THE GET request for student to retrive data from the Database.
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllResults Results
		reqBody, err := ioutil.ReadAll(r.Body)

		// defer the close till after the main function has finished executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllResults)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetResults(db, params["resultsid"], params["studentid"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func timetable(w http.ResponseWriter, r *http.Request) {
	if !pvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn2")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}

	//THE GET request for student to retrive data from the Database.
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllTimetable Timetable
		reqBody, err := ioutil.ReadAll(r.Body)

		// defer the close till after the main function has finished executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllTimetable)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetTimeTable(db, params["timetableid"], params["modulecode"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func commentsrating(w http.ResponseWriter, r *http.Request) {
	if !pvalidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn2")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}
	//THE GET request for student to retrive data from the Database.
	if r.Method == "GET" {
		params := mux.Vars(r)
		var getAllCommentsNRatings CommentsRatings
		reqBody, err := ioutil.ReadAll(r.Body)

		// defer the close till after the main function has finished executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllCommentsNRatings)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("Error in JSON encoding. Error is %s", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information!"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetCommentsRatings(db, params["ratingsid"], params["studentid"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

// func allStudents(w http.ResponseWriter, r *http.Request) {
// 	kv := r.URL.Query()
// 	for k, v := range kv {
// 		fmt.Println(k, v)
// 	}
// 	//returns all the students in JSON
// 	json.NewEncoder(w).Encode(students)
// }

//To check if whether there is a duplicate email in the system
func validateStudentRecord(db *sql.DB, DOB string) bool {
	query := fmt.Sprintf("SELECT * FROM Students WHERE DOB= '%s'", DOB)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var student Students
	for results.Next() {
		err = results.Scan(&student.StudentID, &student.StudentName, &student.DOB, &student.Address, &student.PhoneNumber)
		if err != nil {
			panic(err.Error())
		} else if student.DOB == DOB {
			return true
		}
	}
	return false
}

//Function to validate whether a specific Passenger exists.
func validateStudentID(db *sql.DB, SID string) int {
	query := fmt.Sprintf("SELECT * FROM Students WHERE StudentID=%s", SID)
	var student Students
	row := db.QueryRow(query) //Method to execute the query and is expected to return a single row.
	if err := row.Scan(&student.StudentID, &student.StudentName, &student.DOB, &student.Address, &student.PhoneNumber); err != nil {
		panic(err.Error())
	} else {
		return student.StudentID
	}
}

func validateStudent(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAsgn2")
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if _, err := strconv.Atoi(params["id"]); err != nil { //Converting string to int
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply student information " + "information " + "in JSON format"))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(validateStudentID(db, params["id"])))) //Converting int to string
	}
}

func GetStudentRecord(db *sql.DB, SID string, DOB string) Students {
	results, err := db.Query("SELECT * FROM Students WHERE StudentID=? AND DOB=?", SID, DOB)
	if err != nil {
		panic(err.Error())
	}
	var student Students
	for results.Next() {
		err = results.Scan(&student.StudentID, &student.StudentName, &student.DOB, &student.Address, &student.PhoneNumber)
		if err != nil {
			panic(err.Error())
		}
	}
	return student
}

func EditStudentRecord(db *sql.DB, SID int, SN string, DOB string, ADS string, PN int) bool {
	query := fmt.Sprintf("UPDATE Students SET StudentName='%s', DOB='%s', Address='%s', PhoneNumber=%d WHERE StudentID=%d", SN, DOB, ADS, PN, SID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func DeleteStudents(db *sql.DB, SID int) {
	fmt.Println("Sorry. You are not able to delete your account due to audit purposes.")
}

func GetModules(db *sql.DB, MCD string, SID string) Modules {
	results, err := db.Query("SELECT * FROM Modules WHERE ModuleCode=? AND StudentID=?", MCD, SID)
	if err != nil {
		panic(err.Error())
	}
	var module Modules
	for results.Next() {
		err = results.Scan(&module.ModuleCode, &module.ModuleName, &module.ModuleSynopsis, &module.ModuleObjective, &module.StudentID)
		if err != nil {
			panic(err.Error())
		}
	}
	return module
}

func GetResults(db *sql.DB, RID string, SID string) Results {
	response, err := db.Query("SELECT * FROM Results WHERE ResultsID=? AND StudentID=?", RID, SID)
	if err != nil {
		panic(err.Error())
	}
	var result Results
	for response.Next() {
		err = response.Scan(&result.ResultsID, &result.ResultsGrade, &result.StudentID, &result.ModuleCode)
		if err != nil {
			panic(err.Error())
		}
	}
	return result
}

func GetTimeTable(db *sql.DB, TID string, MCD string) Timetable {
	results, err := db.Query("SELECT * FROM Timetable WHERE TimetableID=? AND ModuleCode=?", TID, MCD)
	if err != nil {
		panic(err.Error())
	}
	var timetable Timetable
	for results.Next() {
		err = results.Scan(&timetable.TimetableID, &timetable.LessonDay, &timetable.StartTime, &timetable.EndTime, &timetable.ModuleCode)
		if err != nil {
			panic(err.Error())
		}
	}
	return timetable
}

func GetCommentsRatings(db *sql.DB, RID string, SID string) CommentsRatings {
	results, err := db.Query("SELECT * FROM CommentsRating WHERE RatingsID=? AND StudentID=?", RID, SID)
	if err != nil {
		panic(err.Error())
	}
	var commentsratings CommentsRatings
	for results.Next() {
		err = results.Scan(&commentsratings.RatingsID, &commentsratings.Ratings, &commentsratings.Comments, &commentsratings.StudentID)
		if err != nil {
			panic(err.Error())
		}
	}
	return commentsratings
}

func main() {
	students = make(map[string]studentInfo)
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/", phome)
	router.HandleFunc("/api/v1/validateStudentRecord/{id}", validateStudent)
	router.HandleFunc("/api/v1/students/{studentid}/{dob}", student).Methods("GET", "PUT", "POST", "DELETE")
	router.HandleFunc("/api/v1/modules/{modulecode}/{studentid}", module).Methods("GET")
	router.HandleFunc("/api/v1/timetable/{timetableid}/{modulecode}", timetable).Methods("GET")
	router.HandleFunc("/api/v1/results/{resultsid}/{studentid}", results).Methods("GET")
	router.HandleFunc("/api/v1/ratings/{ratingsid}/{studentid}", commentsrating).Methods("GET")

	fmt.Println("Listening at port 8103")
	log.Fatal(http.ListenAndServe(":8103", handlers.CORS(headers, methods, origins)(router)))
}
