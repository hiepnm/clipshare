package main

import (
	l4g "code.google.com/p/log4go"
	"net/http"
	"path"
)

func setInteralError(w http.ResponseWriter, err error) {
	if err != nil {
		l4g.Trace(err)
	}
	http.Error(w, "Internal Server Error", 500)
}

func setExistedUser(w http.ResponseWriter, r *http.Request) {
	entry, ok := getEntry(path.Join(config.DocumentRoot, "register/register_error_existeduser.html"))
	if !ok || entry.data == nil {
		setInteralError(w, nil)
		return
	}
	serveFile(w, r, entry)
}