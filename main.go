package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/phyber/negroni-gzip/gzip"
	"net/http"
)

var (
	decoder = schema.NewDecoder()
	client  = &http.Client{}
)

func init() {
	decoder.IgnoreUnknownKeys(true)
}

func MuxHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	herokuWebhookPayload := new(HerokuWebhookPayload)
	err = decoder.Decode(herokuWebhookPayload, r.PostForm)
	if err != nil {
		fmt.Printf("Raw Payload: %v\n", r.PostForm)
		http.Error(w, err.Error(), 422)
		return
	}

	fmt.Printf("Recieved Heroku Deploy Webhook: %v\n", herokuWebhookPayload)

	if NewRelicIsConfigured() {
		go func() {
			handleOutboundRequest("NewRelic", NewRelicRequest(herokuWebhookPayload))
		}()
	}

	if HoneybadgerIsConfigured() {
		go func() {
			handleOutboundRequest("Honeybadger", HoneybadgerRequest(herokuWebhookPayload))
		}()
	}

	if SlackIsConfigured() {
		go func() {
			handleOutboundRequest("Slack", SlackRequest(herokuWebhookPayload))
		}()
	}

	w.WriteHeader(http.StatusAccepted)
}

func handleOutboundRequest(service string, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s %v\n", service, err)
	} else {
		fmt.Printf("OK: %s %v\n", service, resp)
	}
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
