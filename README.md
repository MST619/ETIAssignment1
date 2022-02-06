# ETIAssignment1 - Scroll down for Assignment 2
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
The first point to note is that microservices are only supposed to have one responsibility. Each microservice would be in charge of CRUDing (Creating, Reading, Updating, and Deleting) or GET, POST, PUT and DELETE for their HTTP counterpart for their respective data types from their databases in this example.

E.g. Passenger microservice would only be responsible for GETTING, POSTING, PUTTING, and DELETING passengers. 

Another factor to consider was that microservices must be loosely linked. (to be continued)
Loosely connected architectures (also known as Microservices) are lean, with a single responsibility and few dependencies, allowing teams to work independently, deploy alone, fail independently, and expand independently, resulting in increased business responsiveness.

The minimum requirements of this assignment required us to have at least 2 microservices.
For this particular assignment, there are three microservices.
The overall objective of this assignment was to demonstrate the ability to develop REST APIs
and to make conscientious consideration in designing microservices. 

For this particular assignment, there were three microservices. Each Microservice has its own table, which is responsible for data transfer for each of the object types. The data may then be transmitted and transferred back and forth through the multiple microservices using GET and POST HTTP requests once it has been queried out of the database. The functionality/logic and data management can then be handled within each Microservice, all while adhering to the loosely linked philosophy that Microservices are known for.

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
1. Install the relevant packages needed to run the code, with the exception of the "strconv" package for drivermain.go:
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
2. To run the code, simply click on Run on whichever IDE you're utilising or using your command prompt, type the following:
```
cd ETIAssignment1\Main\Passenger
go run passengermain.go
 
```
```
cd ETIAssignment1\Main\Driver
go run drivermain.go
 
```
```
cd ETIAssignment1\Main\Trip
go run tripmain.go
 
```
3. Once the microservices are running, open up the respective HTML pages. `ETI Ride Sharing.html` for the Passenger, `Driver.html` for driver, and `Trip.HTML` for the Trip. All of which can be found in `\ETIAssignment1\HTML`.




# ETI Assignment 2


# Assignment 2

Assignment was a continuation of what we have done in Assignment 1, just at a much larger scale. From just 2-3 microservices to 23 to work with. Only one would be assigned for you to work on for this assignment. For my particular package, 3.5, it had to do with students. Some of the features of the packages include viewing particulars, updating particulars, viewing the modules taken and etc. For each one of us to utilise our assigned microservices, we would have to use Docker to create images of our microservices, containerize it and then deploy it so that it can be run by others.

# Package 3.5, Student - Features

- 3.5.1. View particulars

![ViewParticulars.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/ViewParticulars.png)

Request URL: 

[](http://localhost:8104/api/v1/students/1/20-10-1985?key=2c78afaf-97da-4816-bbee-9ad239abb296)

- 3.5.2. Update particulars

![UpdateStudentParticulars.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/UpdateStudentParticulars.png)

- 3.5.3. View modules taken

![ModuleHTML.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/ModuleHTML.png)

Request URL: 

[](http://localhost:8104/api/v1/modules/CM/1?key=2c78afaf-97da-4816-bbee-9ad239abb296)

- 3.5.4. View original results

![ResultsHTML.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/ResultsHTML.png)

Request URL: 

[](http://localhost:8104/api/v1/results/0001/1?key=2c78afaf-97da-4816-bbee-9ad239abb296)

- 3.5.6. View timetable

![TimetableHTML.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/TimetableHTML.png)

[](http://localhost:8104/api/v1/timetable/T001/CM?key=2c78afaf-97da-4816-bbee-9ad239abb296)

- 3.5.8. Search for other students

![Search other students.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/GetOther%20students.png)

Request URL: 

[](http://localhost:8104/api/v1/students/2/21-09-1997?key=2c78afaf-97da-4816-bbee-9ad239abb296)

[3.5.9. View other students’ profile, modules, timetable, ratings, and comments](https://www.notion.so/f8f7dfa2d2cc43a49e0acd91193491b7)

# Design considerations of Student Microservice

There would just be one microservice that I would have to work on for assignment 2. Majority of the features that you see above is dependent on other people’s API. These APIs include class, module, timetable etc. The only API that was done by me was to view and edit the particulars.

Since the API was responsible for viewing and editing, it would utilise the GET and PUT method.

```go
//Part of main func in studentmain.go
methods := handlers.AllowedMethods([]string{"GET", "PUT"})
router.HandleFunc("/api/v1/students/{studentid}/", student).Methods("GET", "PUT")

```

# ****Architecture Diagram****

![ETI2 architecture diagram.drawio.png](https://github.com/MST619/ETIAssignment1/blob/main/ETIAsgn2/ETI2%20architecture%20diagram.drawio.png)

As mentioned above, the student package relies on a lot of other people’s API in order for it to function properly. 

# [Nginx](https://www.nginx.com/)

NGINX, stylized as NGINX, nginx or NginX, is a web server that can also be used as a reverse proxy, load balancer, mail proxy and HTTP cache. From my context, it is for serving the HTML files only. Because without the web server, I cannot access the page. How I was accessing the HTML files was that I was going over to the directory where the files were located and opening them manually. This will only allow my computer to view it and not anyone else. The code below is the actual code in my HTML Dockerfile.

```docker
FROM nginx:latest

COPY ./ /usr/share/nginx/html
```

# Setup Instructions

1. Download [GO](https://go.dev/dl/) and [MySQL Community Edition](https://dev.mysql.com/downloads/installer/) database.
2. Launch the MySQL workbench and create a new MySQL connection.
3. Run the following command:

```sql
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost'
```

This will create an account named user with the password 'password'.

1. Run the `ETIAsgn2 database.sql` file, section by section, and step by step.
2. Clone the repository. Install [GitHub desktop](https://desktop.github.com/) and/or follow the steps [here](https://docs.github.com/en/desktop/contributing-and-collaborating-using-github-desktop/adding-and-cloning-repositories/cloning-and-forking-repositories-from-github-desktop).

# Utilising the code

1. Import all the necessary packages needed to run the code.

```go
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
```

## Docker images

Link to docker images:

Student front-end: 

[Docker Hub](https://hub.docker.com/repository/docker/mst619/main_web)

Student API: 

[Docker Hub](https://hub.docker.com/repository/docker/mst619/main_student)

Student DB:

[Docker Hub](https://hub.docker.com/repository/docker/mst619/mysql)

1. Pull the images from docker

```docker
docker pull mst619/studentfinal
docker pull mst619/studenthtmlfinal
docker pull mst619/studentdbfinal
```

2. To run the code, simply click on Run on whichever IDE you're utilising or using your command prompt, type the following: 

```
cd ETIAsgn2\Main\Student
go run studentmain.go
```

3. Once the microservices are running, open up the respective HTML pages. `student.html` for the Student, `timetable.html` for timetable, `module.html` for the module, `results.html` for the results, and `commentsrating.html` for the comments and ratings. All of which can be found in `ETIAssignment1/ETIAsgn2/Main/HTML/`.


