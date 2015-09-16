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

	payload, err := ParseWebhook(r)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}

	if NewRelicIsConfigured() {
		go func() {
			handleOutboundRequest("NewRelic", NewRelicRequest(payload))
		}()
	}

	if HoneybadgerIsConfigured() {
		go func() {
			handleOutboundRequest("Honeybadger", HoneybadgerRequest(payload))
		}()
	}

	if SlackIsConfigured() {
		go func() {
			handleOutboundRequest("Slack", SlackRequest(payload))
		}()
	}

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
