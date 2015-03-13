package main

import (
	"database/sql"
	"net/http"
	"path"
)

func handleActivate(w http.ResponseWriter, r *http.Request) {
	var err error
	con, err := sql.Open("mymysql", config.strDBConnect)
	if err != nil {
		setInteralError(w, err)
		return
	}
	defer con.Close()
	/*check existed uid*/
	uid := r.FormValue("uid")
	value, ok := global.cacheUid.Get(uid)
	if !ok {
		http.NotFound(w, r)
		return
	}
	/*Update to cl_Users table*/
	err = updateTempUserToUserTable(con, uid, value.(string))
	switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			setInteralError(w, err)
			return
	}
	
	entry, ok := getEntry(path.Join(config.DocumentRoot, "register/register_success_active.html"))
	if !ok || entry.data == nil {
		setInteralError(w, nil)
		return
	}
	serveFile(w, r, entry)
}