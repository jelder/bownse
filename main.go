package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"net/http"
)

var (
	client = &http.Client{}
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := ParseWebhook(r)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}

	go func() {
		handleOutboundRequest("NewRelic", NewRelicRequest(state))
	}()
	go func() {
		handleOutboundRequest("Honeybadger", HoneybadgerRequest(state))
	}()
	go func() {
		handleOutboundRequest("Slack", SlackRequest(state))
	}()

	w.WriteHeader(http.StatusAccepted)
}

func handleOutboundRequest(service string, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s %+v\n", service, err)
	} else {
		fmt.Printf("OK: %s %+v\n", service, resp)
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/"+config.Secret, WebhookHandler).Methods("POST")
	http.Handle("/", r)

	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)
	n.Run(config.ListenAddress)
}
