package main

import (
	"database/sql"
	"github.com/nu7hatch/gouuid"
	"time"
)
type TempUser struct {
	userName    string
	password    string
	displayName string
	email       string
	birthDate   time.Time
}
type User struct {
	userName string
	password string
	displayName string
	email string
	birthDate time.Time
}
func checkExistedUserInUsersTable(con *sql.DB, uname string) (bool, error) {
	var username string
	err := con.QueryRow("SELECT username FROM cl_Users WHERE username=?", uname).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func insertUserToTempUserTable(con *sql.DB, temp *TempUser) (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	uid := u4.String()
	_, err = con.Query("INSERT INTO cl_TempUsers VALUES (?, ?, ?, ?, ?, ?)", 
		uid, 
		temp.userName, 
		temp.password, 
		temp.displayName, 
		temp.email, 
		temp.birthDate,
		)
	return uid, err
}

func updateTempUserToUserTable(con *sql.DB, uid string, uname string) error{
	var err error
	tempUser := new(TempUser)
	err = con.QueryRow("SELECT UserName, Password, DisplayName, Email, BirthDate FROM cl_TempUsers WHERE Uuid=?", uid).Scan(&tempUser.userName, &tempUser.password, &tempUser.displayName, &tempUser.email, &tempUser.birthDate)
	if err != nil {
		return err
	}
	_, err = con.Query("DELETE FROM cl_TempUsers WHERE username=?",tempUser.userName)
	if err != nil {
		return err
	}
	
	_, err = con.Query("INSERT INTO cl_Users VALUES(?,?,?,?,?)", 
		tempUser.userName, 
		tempUser.password,
		tempUser.displayName,
		tempUser.email,
		tempUser.birthDate,
		)
	return err
}