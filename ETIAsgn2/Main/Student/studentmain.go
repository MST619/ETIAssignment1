package main

import (
	"bytes"
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
	StudentID   string `json:"StudentID"`
	StudentName string `json:"StudentName"`
	DOB         string `json:"DOB"`
	Address     string `json:"Address"`
	PhoneNumber int    `json:"PhoneNumber"`
}

type Modules struct {
	ModuleCode      int        `json:"ModuleCode"`
	ModuleName      string     `json:"ModuleName"`
	ModuleSynopsis  string     `json:"ModuleSynopsis"`
	ModuleObjective string     `json:"ModuleObjective"`
	ModuleStudent   []Students `json:"ModuleStudent"`
	Results         []Results  `json:"Results"`
}

type Results struct {
	ResultsID    int        `json:"ResultsID"`
	ResultsGrade string     `json:"ResultsGrade"`
	StudentID    []Students `json:"StudentsID"`
	ModuleCode   []Modules  `json:"ModuleCode"`
}

type Timetable struct {
	TimetableID string  `json:"TimetableID"`
	ModuleCode  int     `json:"ModuleCode"`
	LessonDay   string  `json:"lesson_day"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Module      Modules `json:"module"`
}

type Ratings struct {
	RatingsID string `json:"RatingsID"`
	StudentID int    `json:"StudentID"`
	Ratings   int    `json:"Ratings"`
	Comments  string `json:"Comments"`
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
	var stdnt Students
	params := mux.Vars(r)
	SID := params["StudentID"]
	studentID, err := strconv.Atoi(SID)
	if studentID == 0 || !validateStudentID(SID) || err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply student information " + "information " + "in JSON format"))
		return
	}
	//std = GetStudentRecord(SID)

	// db, err := sql.Open("mysql", "user:password@tcp(studentdb:3306)/ETIAsgn2")
	// if err != nil {
	// 	panic(err.Error())
	// } else {
	// 	fmt.Println("Database opened!")
	// }

	//THE GET request for student to retrive data from the Database.
	if r.Method == "GET" {
		stdnt = GetStudentRecord(studentID, stdnt.DOB)
		if stdnt == (Students{}) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("No student found"))
		} else {
			json.NewEncoder(w).Encode(stdnt)
			w.WriteHeader(http.StatusAccepted)
		}
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
		json.NewEncoder(w).Encode(studentID)
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "PUT" {
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				println(err.Error())
			}
			defer r.Body.Close()
			var updateStudent Students
			err = json.Unmarshal(reqBody, &updateStudent)
			if !EditStudentRecord(updateStudent) || err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("No student found with: " + updateStudent.DOB))
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("201 - Student updated!"))
				return
			}
		}

	} else {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply student information " + "information " + "in JSON format"))
		return
	}
}

func getmodule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	method := params["method"]
	SIDParam := params["StudentID"]
	SID, err := strconv.Atoi(SIDParam)

	if method == "" || err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Please supply student's information and valid method"))
		return
	} else {
		switch string(method) {
		case "getModules":
			Modules := getModulesTaken(SID)
			if len(Modules) == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("No modules found"))
			} else {
				json.NewEncoder(w).Encode(Modules)
			}
		case "getResults":
			Results := getResults(SID)
			if len(Results) == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("No results found"))
			} else {
				json.NewEncoder(w).Encode(Results)
				w.WriteHeader(http.StatusAccepted)
			}
		case "getTimetable":
			Timetable := getTimeTable(SID)
			println(Timetable)
			w.WriteHeader(http.StatusAccepted)

		case "getAdjustedResults":
			Student := getAdjustedResults(SID)
			if len(Student) == 0 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Adjusted results are empty"))
			} else {
				json.NewEncoder(w).Encode(Student)
				w.WriteHeader(http.StatusAccepted)
			}
		}
	}
}

func otherdetails(w http.ResponseWriter, r *http.Request) {
	var studentList []Students
	URL := "http://localhost:8103/api/v1/getAllStudentsWithRatings"
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
	}
	if response.StatusCode == http.StatusAccepted {
		results, err := ioutil.ReadAll(response.Body)
		if err != nil || json.Unmarshal([]byte(results), &studentList) != nil {
			println(err)
		}
	}
	json.NewEncoder(w).Encode(studentList)
}

func allStudents(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, v)
	}
	//returns all the students in JSON
	json.NewEncoder(w).Encode(students)
}

//To check if whether there is a duplicate email in the system
func validateStudentRecord(DOB string) bool {
	URL := fmt.Sprintf("http://localhost:8103/api/v1/validateStudentRecord/%s", DOB)
	//query := fmt.Sprintf("SELECT * FROM Students WHERE DOB= '%s'", DOB)
	response, err := http.Get(URL)
	if err != nil {
		panic(err.Error())
	}
	if response.StatusCode == 202 {
		response, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var std Students
			err = json.Unmarshal([]byte(response), &std)
			if err != nil {
				return true
			}
		}
	}
	return false
}

//Function to validate whether a specific student exists.
func validateStudentID(SID string) bool {
	URL := fmt.Sprintf("http://localhost:8103/api/v1/validateStudentID/%s", SID)
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
	}
	if response.StatusCode == 202 {
		response, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var std Students
			err = json.Unmarshal([]byte(response), &std)
			if err != nil {
				return true
			}
		}
	}
	return false
}

//3.5.1.	View particulars
func GetStudentRecord(SID int, DOB string) Students {
	URL := fmt.Sprintf("http://localhost:8103/api/v1/GetStudentRecord/%d", SID)
	response, err := http.Get(URL)
	var student Students
	if err != nil {
		panic(err.Error())
	}
	if response.StatusCode == http.StatusAccepted {
		response, err := ioutil.ReadAll(response.Body)
		if err == nil {
			err = json.Unmarshal(response, &student)
		}
		println(err)
	}
	return student
}

//3.5.3.	View modules taken
func getModulesTaken(SID int) []Modules {
	// results, err := db.Query("SELECT * FROM Modules WHERE StudentID=?", SID)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// var modules Modules
	// for results.Next() {
	// 	err = results.Scan(&modules.ModuleCode, &modules.ModuleName, &modules.ModuleSynopsis, &modules.ModuleObjective, &modules.ModuleStudent)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// }
	// return modules

	//FOR WHEN I GET THE ENDPOINT FROM PKG3.4
	URL := fmt.Sprintf("http://localhost:8103/api/v1/getModulesTaken/%d", SID)
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
	} else if response.StatusCode == http.StatusAccepted {
		reqBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var mods []Modules
			err := json.Unmarshal(reqBody, &mods)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	return nil
}

//3.5.4.	View original results
func getResults(SID int) []Results {
	URL := "http://localhost:8103/api/v1/getResults"
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
	} else if response.StatusCode == http.StatusAccepted {
		respose, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var ResultsData []Results
			err := json.Unmarshal(respose, &ResultsData)
			if err != nil {
				panic(err.Error())
			}
			return ResultsData
		}
	}
	return nil
}

//3.5.5.	View adjusted results after marks trading
func getAdjustedResults(SID int) []Results {
	URL := "http://localhost:8103/api/v1/getResults"
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(err.Error())
	} else if response.StatusCode == http.StatusAccepted {
		response, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var adjResults []Results
			err := json.Unmarshal(response, &adjResults)
			if err != nil {
				panic(err.Error())
			}
			return adjResults
		}
	}
	return nil
}

//3.5.6.	View timetable
func getTimeTable(SID int) bool {
	URL := fmt.Sprintf("http://localhost:8103/api/v1/getTimetable/%d", SID)
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
	} else if response.StatusCode == http.StatusAccepted {
		reqBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var timetable Timetable
			err := json.Unmarshal(reqBody, &timetable)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	return false
}

//3.5.7.	List all students with ratings
func getAllStudentsWithRatings() []Ratings {
	response, err := http.Get("http://localhost:8103/api/v1/getAllStudentsWithRatings")
	if err != nil {
		fmt.Print(err.Error())
	}
	if response.StatusCode == http.StatusAccepted {
		reqBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			var RCs []Ratings
			err := json.Unmarshal(reqBody, &RCs)
			if err == nil {
				return RCs
			}
		}
	}
	return nil
}

//3.5.8.	Search for other students
func getDiffStudent(SID string) Students {
	url := "http://localhost:8103/api/v1/getDiffStudent/1"
	reqBody, err := http.Get(url)
	var student Students

	if err != nil {
		fmt.Print(err.Error())
	}
	if reqBody.StatusCode == http.StatusAccepted {
		reqBody, err := ioutil.ReadAll(reqBody.Body)
		if err != nil || json.Unmarshal([]byte(reqBody), &student) != nil {
			print(err)
		}
	}
	return student
}

//3.5.2.	Update particulars
func EditStudentRecord(student Students) bool {
	json, _ := json.Marshal(student)
	URL := "http://localhost:8103/api/v1/EditStudentRecord"

	response, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(json))
	if err != nil {
		panic(err.Error())
	} else {
		response.Header.Set("Content-Type", "application/json")

		request := &http.Client{}
		response, err := request.Do(response)

		if err != nil {
			fmt.Printf("Error in JSON encoding. Error is %s", err)
		} else {
			if response.StatusCode == http.StatusCreated {
				response.Body.Close()
			}
		}
		response.Body.Close()
	}
	return false
}

func DeleteStudents(db *sql.DB, SID int) {
	fmt.Println("Sorry. You are not able to delete your account due to audit purposes.")
}

func main() {

	students = make(map[string]studentInfo)
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/", phome)
	//router.HandleFunc("/api/v1/validateStudentRecord/{id}", validateStudent)
	router.HandleFunc("/api/v1/students/{studentid}/", student).Methods("GET", "PUT")
	router.HandleFunc("/api/v1/getmodule/{method}/{studentid}", getmodule).Methods("GET")
	router.HandleFunc("/api/v1/otherdetails/{method}/{studentname}", otherdetails).Methods("GET")

	fmt.Println("Listening at port 8103")
	log.Fatal(http.ListenAndServe(":8103", handlers.CORS(headers, methods, origins)(router)))
}
