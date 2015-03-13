package main

import (
	"bytes"
	l4g "code.google.com/p/log4go"
	"compress/gzip"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
	"io"
	"os"
)

type FileHandler struct {
	
}
type cacheEntry struct {
	data          []byte
	checksum      string
	compressData []byte
	time.Time
}
func checkModifiedSince(w http.ResponseWriter, r *http.Request, e *cacheEntry) bool {
	t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))
	if err == nil && e.Before(t.Add(1*time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	return false
}

func checkEtag(w http.ResponseWriter, r *http.Request, e *cacheEntry) bool {
	etag := r.Header.Get("If-None-Match")
	if etag != "" {
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	return false
}
func readFileSystem(rpath string) (reply interface{}, err error) {
	fpath := path.Join(config.DocumentRoot, rpath)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	finfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	
	buffCompress := &bytes.Buffer{}
	z := gzip.NewWriter(buffCompress)
	defer z.Close()
	
	e := &cacheEntry{data: nil, checksum: "", compressData: nil, Time: finfo.ModTime()}
	e.data = make([]byte, 0)
	for {
		buf := make([]byte, 1<<16)//64KB
		n, err := file.Read(buf)
		
		if err == io.EOF {
			e.data = append(e.data, buf[:n]...)
			break
		}
		z.Write(buf)
		e.data = append(e.data, buf[:n]...)
	}
	compressData := buffCompress.Bytes()
	e.compressData = make([]byte, len(compressData))
	e.compressData = compressData
	reply = e
	return
}

func getEntry(path string) (*cacheEntry, bool) {
	val, ok := global.cacheFile.Get(path)//get data from local cache
	if ok {
		return val.(*cacheEntry), ok
	}
	
	data, err := readFileSystem(path)//get data from DB or filesystem
	if err == nil {
		e := data.(*cacheEntry)
		global.cacheFile.Add(path, e)
		return e, true
	} else {
		l4g.Trace(err)
	}
	return nil, false
}

func addHeader(w http.ResponseWriter, r *http.Request) {
	/*server tag*/
	if len(config.serverTag) > 0 {
		w.Header().Set("Server", config.serverTag)
	}
	/*cache control*/
	w.Header().Set("Cache-Control", "max-age=600, must-revalidate, proxy-revalidate")
	t := time.Now()
	w.Header().Set("Expires", t.Add(600*time.Second).UTC().Format(http.TimeFormat))
	/*access*/
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mm := mime.TypeByExtension(path.Ext(r.URL.Path))
	w.Header().Set("Content-Type", mm)
}

func response(w http.ResponseWriter, r *http.Request, entry *cacheEntry) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Write(entry.data)
	} else {
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(entry.compressData)
	}
}
func serveFile(w http.ResponseWriter, r *http.Request, entry *cacheEntry) {
	addHeader(w, r)
	if checkModifiedSince(w, r, entry) || checkEtag(w, r, entry) {
		return
	}
	response(w, r, entry)
	
}
func (f *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entry, ok := getEntry(r.URL.Path)
	if !ok || entry.data == nil {
		if r.URL.Path == "/" {
			w.Write([]byte("OK"))
		} else {
			http.NotFound(w, r)
		}
	}
	serveFile(w, r, entry)
}
