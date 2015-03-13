package main

import (
	l4g "code.google.com/p/log4go"
	"net/http"
	"log"
	"runtime"
)
var (
	global *Global
	config *ConfigGlobal
)
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initialize()
	l4g.Trace("ClipShare Server now listen at: %v", config.bindAddress)
	registerHandler()
	go global.fileServer.ListenAndServe()
	if err:=http.ListenAndServe(config.controlAddress, nil); err!=nil {
		log.Panicln(err)
	}
}