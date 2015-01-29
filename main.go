package main

import (
	"bytes"
	// "fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func MuxHandler(w http.ResponseWriter, r *http.Request) {
	if len(TargetUrls) < 1 {
		http.NotFound(w, r)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	for _, url := range TargetUrls {
		go Repost(url, body)
	}
	w.WriteHeader(201)
}

func Repost(url string, payload []byte) {
	log.Println(url, len(payload), "bytes")
	client.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(payload))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := mux.NewRouter()
	r.HandleFunc("/"+Secret, MuxHandler).Methods("POST")
	http.Handle("/", r)

	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
	n.UseHandler(r)
	n.Run(ListenAddress)
}
