#/-------------------SECTION 1---------------------\
#1. CREATE THE DATABASE
CREATE DATABASE ETIAsgn2;

#2. CREATE THE PASSENGER TABLE
USE ETIAsgn2;

#3. CREATE AND GRANT PERMISSION TO THE USER
CREATE USER 'user'@'%' IDENTIFIED BY 'password';

GRANT ALL PRIVILEGES
ON ETIAsgn2.*
TO 'user'@'%';

#/-------------------SECTION 2---------------------\
#1. CREATE THE STUDENT TABLE
CREATE TABLE Students 
(
	StudentID VARCHAR(5) NOT NULL PRIMARY KEY, 
	StudentName VARCHAR(30),
	DOB VARCHAR(10), 
	Address VARCHAR(50), 
	PhoneNumber VARCHAR(8)
);

#2. INSERT VALUES INTO THE NEWLY CREATED STUDENTS TABLE
INSERT INTO Students (StudentID, StudentName, DOB, Address, PhoneNumber) VALUES ("1", "Jake", "20-10-1985", "1 Temasek Avenue", "99998888");
INSERT INTO Students (StudentID, StudentName, DOB, Address, PhoneNumber) VALUES ("2", "Amy", "21-09-1997", "1057 Eunos Avenue", "88889999");
INSERT INTO Students (StudentID, StudentName, DOB, Address, PhoneNumber) VALUES ("3", "Raymond", "15-12-1995", "1 Bilal Lane", "77776666");
INSERT INTO Students (StudentID, StudentName, DOB, Address, PhoneNumber) VALUES ("4", "Terry", "14-02-1976", "Kallang Puddin Rd", "66665555");
INSERT INTO Students (StudentID, StudentName, DOB, Address, PhoneNumber) VALUES ("5", "Charles", "09-08-1965", "163 Tanglin Rd", "55554444");

#3. CHECK TO SEE IF ALL THE DATA HAS BEEN ADDED
SELECT * FROM Students;