#/-------------------SECTION 1---------------------\
#1. CREATE THE DATABASE
CREATE DATABASE ETIAsgn;

#2. CREATE THE PASSENGER TABLE
USE ETIAsgn;


#/-------------------SECTION 2---------------------\
#1. CREATE THE PASSENGER TABLE
CREATE TABLE Passengers (
PassengerID VARCHAR(5) NOT NULL PRIMARY KEY, 
FirstName VARCHAR(30), 
LastName VARCHAR(30), 
PhoneNumber VARCHAR(8), 
Email VARCHAR(50)
);

#2. INSERT VALUES INTO THE NEWLY CREATED PASSENGERS TABLE
INSERT INTO Passengers (PassengerID, FirstName, LastName, PhoneNumber, Email) VALUES ("1", "Jake", "Peralta", "99998888", "noice@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, PhoneNumber, Email) VALUES ("2", "Amy", "Santiago", "88887777", "asb99@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, PhoneNumber, Email) VALUES ("3", "Raymond", "Holt", "77776666", "rayholt99@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, PhoneNumber, Email) VALUES ("4", "Terry", "Jeffords", "66665555", "yoghurt@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, PhoneNumber, Email) VALUES ("5", "Charles", "Boyle", "55554444", "cbb99@gmail.com");

#3. CHECK TO SEE IF ALL THE DATA HAS BEEN ADDED
SELECT * FROM Passengers;



#/-------------------SECTION 3---------------------\
#1. CREATE THE DRIVERS TABLE
CREATE TABLE Drivers (
DriverID VARCHAR(5) NOT NULL PRIMARY KEY, 
FirstName VARCHAR(30), 
LastName VARCHAR(30), 
PhoneNumber VARCHAR(8), 
Email VARCHAR(50), 
LicenseNo VARCHAR(6)
);

#2. INSERT VALUES INTO THE NEWLY CREATED DRIVERS TABLE
INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ("1", "Max", "Verstappen", "98765432", "mv33@gmail.com", "999888");
INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ("2", "Lewis", "Hamilton", "12345678", "lh44@gmail.com", "888777");
INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ("3", "Sergio", "Perez", "89621515", "sp11@gmail.com", "777666");
INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ("4", "Lando", "Norris", "29958585", "ln4@gmail.com", "666555");
INSERT INTO Drivers (DriverID, FirstName, LastName, PhoneNumber, Email, LicenseNo) VALUES ("5", "Charles", "Leclerc", "55317916", "cl16@gmail.com", "555444");

#3. CHECK TO SEE IF ALL THE DATA HAS BEEN ADDED
SELECT * FROM Drivers;



#/-------------------SECTION 4---------------------\
#1. CREATE THE TRIPS TABLE
CREATE TABLE Trips (
	TripID VARCHAR(5) NOT NULL, 
    PickupPC VARCHAR(10) NOT NULL, 
    DropoffPC VARCHAR(10) NOT NULL, 
    DriverID VARCHAR(5), 
    PassengerID VARCHAR(5), 
    TripStatus VARCHAR(15) NOT NULL,
    PRIMARY KEY (TripID),
    FOREIGN KEY (DriverID) REFERENCES Drivers(DriverID),
    FOREIGN KEY (PassengerID) REFERENCES Passengers(PassengerID)
);

#2. INSERT VALUES INTO THE NEWLY CREATED DRIVERS TABLE
INSERT INTO Trips (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ("1", "237994", "218031", "1", "5", "Finished");
INSERT INTO Trips (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ("2", "738733", "608586", "2", "4", "Processing");
INSERT INTO Trips (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ("3", "329852", "199597", "3", "3", "In Progress");
INSERT INTO Trips (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ("4", "128384", "058455", "4", "2", "Finished");
INSERT INTO Trips (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ("5", "369974", "428751", "5", "1", "Finished");

#3. CHECK TO SEE IF ALL THE DATA HAS BEEN ADDED
SELECT * FROM Trips
