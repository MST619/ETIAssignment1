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


## Setup Instructions 
1. Download [GO](https://go.dev/dl/) and [MySQL Community Edition](https://dev.mysql.com/downloads/installer/) database.
2. Launch the MySQL workbench and create a new MySQL connection. 
3. Run the following command:
```
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost'
 
```
This will create an account named user with the password 'password'.

4. Run the `ETIAsgn database.sql` file, section by section.




