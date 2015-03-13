package main

import (
	"database/sql"
	"net/http"
)

func handleCheckExistedUser(w http.ResponseWriter, r *http.Request) {
	var err error
	con, err := sql.Open("mymysql", config.strDBConnect)
	if err != nil {
		setInteralError(w, err)
		return
	}
	username := r.FormValue("username")
	existed, err := checkExistedUserInUsersTable(con, username)
	if err != nil {
		setInteralError(w, err)
		return
	}
	var ret string
	if existed {
		ret = "1"
	} else {
		ret = "0"
	}
	w.Write([]byte(ret))
}