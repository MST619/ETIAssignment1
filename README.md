# ETIAssignment1
Student ID: S10197943E  
Student Name: Min Se Thu

## The Assignment
The project is one of two assignments that we will be involved in, for our ETI module.
The assignment is built with the following:
1. Golang (GO),
2. HTML,
3. JavaScript,
4. jQuery.Ajax, and
5. MySQL database

## Design Considerations
The minimum requirements of this assignment required us to have at least 2 microservices.
For this particular assignment, there are three microservices.
The overall objective of this assignment was to demonstrate the ability to develop REST APIs
and to make conscientious consideration in designing microservices. 
For this particular assignment, there were three microservices: 

1. Passenger  
The passenger microservice makes use of the POST, GET, and PUT HTTP method.

2. Driver  
The driver microservice also makes use of the POST, GET and PUT HTTP method.
 
3. Trip  
The trip microservice has it's own POST, GET and PUT method but also makes use of
the Passenger and Driver microservices by calling them. 

## Archicture Diagram
Here is an architecture diagram to visualise how the application works and the relationship between the frontend, the microservices, and the database.  
![ETI Ride Sharing Architecture diagram](https://raw.githubusercontent.com/MST619/ETIAssignment1/main/ETI%20architecture%20diagram2.png)

## Setup Instructions 
1. Download [GO](https://go.dev/dl/) and [MySQL Community Edition](https://dev.mysql.com/downloads/installer/) database.
2. Launch the MySQL workbench and create a new MySQL connection. 
3. Run the following command:
```
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost'
 
```
This will create an account named user with the password 'password'.

4. Run the `ETIAsgn database.sql` file, section by section, ans step by step.

5. Clone the repository. Install [GitHub desktop](https://desktop.github.com/) and/or follow the steps [here](https://docs.github.com/en/desktop/contributing-and-collaborating-using-github-desktop/adding-and-cloning-repositories/cloning-and-forking-repositories-from-github-desktop)

## Utilising the code
6. Install the relevant packages needed to run the code, with the exception of the "strconv" pacakge for drivermain.go:
```
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
 
```




