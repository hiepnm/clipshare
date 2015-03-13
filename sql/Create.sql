CREATE DATABASE `ClipShare`;
GRANT ALL PRIVILEGES ON `ClipShare`.* TO "hiepnm"@"localhost" IDENTIFIED BY "123456";
FLUSH PRIVILEGES;
USE `ClipShare`;
CREATE TABLE `cl_Users` (
	`UserName` VARCHAR(20) NOT NULL,
	`PassWord` VARCHAR(30) NOT NULL,
	`DisplayName` VARCHAR(40) NOT NULL,
	`Email` VARCHAR(40) NOT NULL,
	`BirthDate` DATETIME NOT NULL,
	PRIMARY KEY (`UserName`)
);
CREATE TABLE `cl_TempUsers` (
	`Uuid` VARCHAR(36) NOT NULL,
	`UserName` VARCHAR(20) NOT NULL,
	`PassWord` VARCHAR(30) NOT NULL,
	`DisplayName` VARCHAR(40) NOT NULL,
	`Email` VARCHAR(40) NOT NULL,
	`BirthDate` DATETIME NOT NULL,
	PRIMARY KEY (`Uuid`)
);

CREATE TABLE `cl_Clips` (
	`ClipID` VARCHAR(36) NOT NULL,
	`ClipName` VARCHAR(20) NOT NULL,
	`Owner` VARCHAR(30) NOT NULL,
	`UpDate` DATETIME NOT NULL,
	PRIMARY KEY (`ClipID`)
);