package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/phyber/negroni-gzip/gzip"
	// "log"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	decoder = schema.NewDecoder()
)

func MuxHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	herokuWebhookPayload := new(HerokuWebhookPayload)
	err = decoder.Decode(herokuWebhookPayload, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Recieved Heroky Deploy Webhook: %v\n", herokuWebhookPayload)
	if NewRelicIsConfigured() {
		newrelic <- herokuWebhookPayload
	}
	if HoneybadgerIsConfigured() {
		honeybadger <- herokuWebhookPayload
	}

	w.WriteHeader(http.StatusAccepted)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/"+config.Secret, MuxHandler).Methods("POST")
	http.Handle("/", r)

	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)
	n.Run(config.ListenAddress)
}
