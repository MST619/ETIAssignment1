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


#/-------------------SECTION 3---------------------\
#1. CREATE THE STUDENT TABLE
CREATE TABLE Modules
(
	ModuleCode VARCHAR(5) NOT NULL PRIMARY KEY, 
	ModuleName VARCHAR(50),
	ModuleSynopsis VARCHAR(50), 
	ModuleObjective VARCHAR(50),
    StudentID VARCHAR(5),
    FOREIGN KEY (StudentID) REFERENCES Students(StudentID)
);

INSERT INTO Modules (ModuleCode, ModuleName, ModuleSynopsis, ModuleObjective, StudentID) VALUES ("CM", "Computing Maths", "Maths for computing", "Instill CM in students", "1");
INSERT INTO Modules (ModuleCode, ModuleName, ModuleSynopsis, ModuleObjective, StudentID) VALUES ("CSF", "Cyber Security Fundamentals", "Tech students about CS fundamentals", "Instill CSF in students", "2");
INSERT INTO Modules (ModuleCode, ModuleName, ModuleSynopsis, ModuleObjective, StudentID) VALUES ("DP", "Data Protection", "Protect students data", "Instill DP in students", "3");
INSERT INTO Modules (ModuleCode, ModuleName, ModuleSynopsis, ModuleObjective, StudentID) VALUES ("PRG1", "Programming 1", "Python Programming for students", "Instill PRG1 in students", "4");
INSERT INTO Modules (ModuleCode, ModuleName, ModuleSynopsis, ModuleObjective, StudentID) VALUES ("DB", "Databases", "Backend databases for students", "Instill DB in students", "5");

SELECT * FROM Modules;


#/-------------------SECTION 4---------------------\
CREATE TABLE Results
(
	ResultsID VARCHAR(5) NOT NULL PRIMARY KEY,
    ResultsGrade VARCHAR(5),
    StudentID VARCHAR(5),
    ModuleCode VARCHAR(5),
    FOREIGN KEY (StudentID) REFERENCES Students(StudentID),
    FOREIGN KEY (ModuleCode) REFERENCES Modules(ModuleCode)
);

INSERT INTO Results (ResultsID, ResultsGrade, StudentID, ModuleCode) VALUES ("0001", "A", "1", "CM");
INSERT INTO Results (ResultsID, ResultsGrade, StudentID, ModuleCode) VALUES ("0002", "B", "3", "CSF");
INSERT INTO Results (ResultsID, ResultsGrade, StudentID, ModuleCode) VALUES ("0003", "C", "5", "DB");
INSERT INTO Results (ResultsID, ResultsGrade, StudentID, ModuleCode) VALUES ("0004", "D", "4", "PRG1");
INSERT INTO Results (ResultsID, ResultsGrade, StudentID, ModuleCode) VALUES ("0005", "F", "2", "DP");

SELECT * FROM Results;


#/-------------------SECTION 5---------------------\
CREATE TABLE Timetable
(
	TimetableID VARCHAR(5) NOT NULL PRIMARY KEY,
    LessonDay VARCHAR(10),
    StartTime VARCHAR(5),
    EndTime VARCHAR(5),
    ModuleCode VARCHAR(5),
    FOREIGN KEY (ModuleCode) REFERENCES Modules(ModuleCode)
);

INSERT INTO Timetable (TimetableID, LessonDay, StartTime, EndTime, ModuleCode) VALUES ("T001", "Monday", "9AM", "11AM", "CM");
INSERT INTO Timetable (TimetableID, LessonDay, StartTime, EndTime, ModuleCode) VALUES ("T002", "Tuesday", "10AM", "12PM", "CSF");
INSERT INTO Timetable (TimetableID, LessonDay, StartTime, EndTime, ModuleCode) VALUES ("T003", "Wednesday", "11AM", "1PM", "DB");
INSERT INTO Timetable (TimetableID, LessonDay, StartTime, EndTime, ModuleCode) VALUES ("T004", "Thursday", "1PM", "3PM", "DP");
INSERT INTO Timetable (TimetableID, LessonDay, StartTime, EndTime, ModuleCode) VALUES ("T005", "Friday", "4PM", "6PM", "PRG1");

SELECT * FROM Timetable;


CREATE TABLE CommentsRating 
(
	RatingsID VARCHAR(5) NOT NULL PRIMARY KEY,
    Ratings VARCHAR(3),
    Comments VARCHAR(50),
    StudentID VARCHAR(5),
    FOREIGN KEY (StudentID) REFERENCES Students(StudentID)
);

INSERT INTO CommentsRating (RatingsID, Ratings, Comments, StudentID) VALUES ("CR01", "5*", "Hard but very rewarding module", "1");
INSERT INTO CommentsRating (RatingsID, Ratings, Comments, StudentID) VALUES ("CR02", "4*", "Fun and interesting module", "3");
INSERT INTO CommentsRating (RatingsID, Ratings, Comments, StudentID) VALUES ("CR03", "3*", "Okay module", "5");
INSERT INTO CommentsRating (RatingsID, Ratings, Comments, StudentID) VALUES ("CR04", "2*", "Just passed programming", "4");
INSERT INTO CommentsRating (RatingsID, Ratings, Comments, StudentID) VALUES ("CR05", "1*", "Failed module", "2");

SELECT * FROM CommentsRating;