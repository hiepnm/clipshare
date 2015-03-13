package main

import (
	"github.com/hashicorp/golang-lru"
	l4g "code.google.com/p/log4go"
	"flag"
	"github.com/vaughan0/go-ini"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"strconv"
)

type Global struct {
	fileServer http.Server
	cacheUid *lru.Cache
	cacheFile *lru.Cache
}
type ConfigGlobal struct {
	controlAddress    string
	controlAddressSSL string
	bindAddress       string
	bindAddressSSL    string
	strDBConnect      string
	serverTag       string
	ServerMail      string
	activeEmail     string
	activeEmailPass string
	templateEmail   string
	DocumentRoot    string
	DEVELOPMENT     bool
	CacheUIDSize int64
	CacheFileSize int64
}

func initialize() {
	global = new(Global)
	config = &ConfigGlobal{
		controlAddress:    "127.0.0.1:44401",
		controlAddressSSL: "127.0.0.1:44434",
		bindAddress:       "127.0.0.1:44400",
		bindAddressSSL:    "127.0.0.1:44433",
		serverTag:         "cs/366",
		DocumentRoot:      "/var/html",
		DEVELOPMENT:       false,
	}
	var conf = flag.String("f", "", "Config file")
	var logc = flag.String("l", "", "log config file")
	var developmentConf = "0"
	flag.Parse()
	if len(*logc) > 0 {
		l4g.LoadConfiguration(*logc)
	}
	if len(*conf) > 0 {
		confFile, e := ini.LoadFile(*conf)
		if e != nil {
			log.Panicln(e)
		}
		if bindAddr, ok := confFile.Get("", "bindaddress"); ok {
			config.bindAddress = bindAddr
		}
		if bindAddrSSL, ok := confFile.Get("", "bindaddressSSL"); ok {
			config.bindAddressSSL = bindAddrSSL
		}
		if controlAddr, ok := confFile.Get("", "controladdress"); ok {
			config.controlAddress = controlAddr
		}
		if controlAddrSSL, ok := confFile.Get("", "controladdressSSL"); ok {
			config.controlAddressSSL = controlAddrSSL
		}
		if stag, ok := confFile.Get("", "servertag"); ok {
			config.serverTag = stag
		}
		if dev, ok := confFile.Get("", "DEVELOPMENT"); ok {
			developmentConf = dev
		}
		if serverMail, ok := confFile.Get("", "servermail"); ok {
			config.ServerMail = serverMail
		}
		if email, ok := confFile.Get("", "activeemail"); ok {
			config.activeEmail = email
		}
		if emailpass, ok := confFile.Get("", "activeemailpass"); ok {
			config.activeEmailPass = emailpass
		}
		if docroot, ok := confFile.Get("", "activeemailpass"); ok {
			config.DocumentRoot = docroot
		}
		if fileTemplate, ok := confFile.Get("", "templatefile"); ok {
			data, err := ioutil.ReadFile(path.Join(config.DocumentRoot, fileTemplate))
			if err != nil {
				log.Panicln(err)
			}
			config.templateEmail = string(data)
		}
		var dbConfig []string = make([]string, 3)
		config.strDBConnect = ""
		if dbname_, ok := confFile.Get("", "dbname"); ok {
			dbConfig[0] = dbname_
		}
		if admin_, ok := confFile.Get("", "admin"); ok {
			dbConfig[1] = admin_
		}
		if dbpass_, ok := confFile.Get("", "dbpass"); ok {
			dbConfig[2] += dbpass_
		}
		if len(dbConfig) != 3 {
			log.Panicln("no database infomation")
		}
		config.strDBConnect = strings.Join(dbConfig, "/")
		
		if cacheUidSize, ok := confFile.Get("", "CacheUIDSize"); ok {
			nCacheUidSize, err := strconv.ParseInt(cacheUidSize, 10, 0)
			if err != nil {
				log.Panicln(err)
			}
			config.CacheUIDSize = nCacheUidSize
			global.cacheUid, err = lru.New(int(nCacheUidSize))
			if err != nil {
				log.Panic(err)
			}
		}
		
		if cacheFileSize, ok := confFile.Get("", "CacheFileSize"); ok {
			nCacheFileSize, err := strconv.ParseInt(cacheFileSize, 10, 0)
			if err != nil {
				log.Panicln(err)
			}
			config.CacheFileSize = nCacheFileSize
			global.cacheFile, err = lru.New(int(nCacheFileSize))
			if err != nil {
				log.Panic(err)
			}
		}
	}
	config.DEVELOPMENT = (os.Getenv("DEVELOPMENT") == "1" || developmentConf == "1")
}

func registerHandler() {
	/*control handler*/
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/active", handleActivate)
	http.HandleFunc("/check_existed_user", handleCheckExistedUser)
	/*file serving*/
	fileHandler := &FileHandler{}
	global.fileServer = http.Server{Handler: fileHandler, Addr: config.bindAddress, ReadTimeout: 10 * time.Second}
}
