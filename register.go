package main

import (
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
	"net/http"
	"time"
	"path"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var err error
	con, err := sql.Open("mymysql", config.strDBConnect)
	if err != nil {
		setInteralError(w, err)
		return
	}
	defer con.Close()
	/*check existed*/
	uname := r.FormValue("username")
	email := r.FormValue("email")
	existed, err := checkExistedUserInUsersTable(con, uname)
	if err != nil {
		setInteralError(w, err)
		return
	}
	if existed {
		setExistedUser(w, r)
		return
	}

	/*Insert into TempUser table*/
	birthDate, err := time.Parse("29/09/1988", r.FormValue("BirthDate"))
	if err != nil {
		setInteralError(w, err)
		return
	}
	tempUser := &TempUser{
		displayName: r.FormValue("DisplayName"),
		userName:    uname,
		password:    r.FormValue("Password"),
		email:       email,
		birthDate:   birthDate,
	}
	uid, err := insertUserToTempUserTable(con, tempUser)
	if err != nil {
		setInteralError(w, err)
		return
	}
	global.cacheUid.Add(uid, tempUser.userName)
	/*send email*/
	err = sendActivateTempUserMail(email, uid, tempUser)
	if err != nil {
		setInteralError(w, err)
		return
	}
	
	entry, ok := getEntry(path.Join(config.DocumentRoot, "register/register_wait_active.html"))
	if !ok || entry.data == nil {
		setInteralError(w, nil)
		return
	}
	serveFile(w, r, entry)
}