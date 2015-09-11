// http://docs.honeybadger.io/article/38-deployment-tracking-on-heroku
// https://api.honeybadger.io/v1/deploys?deploy[environment]=production&deploy[local_username]={{user}}&deploy[revision]={{head}}&api_key=asdf

package main

import (
	"fmt"
	"net/http"
	"net/url"
)

var (
	honeybadger = make(chan *HerokuWebhookPayload)
)

func init() {
	if HoneybadgerIsConfigured() {
		go handleHoneybadger()
	} else {
		fmt.Println("Honeybadger is not full configured")
	}
}

func HoneybadgerIsConfigured() bool {
	return ENV["HONEYBADGER_API_KEY"] != ""
}

func handleHoneybadger() {
	for {
		payload := <-honeybadger
		params := url.Values{
			"deploy[environment]":    {payload.Environment()},
			"deploy[local_username]": {payload.User},
			"deploy[revision]":       {payload.Head},
			"api_key":                {ENV["HONEYBADGER_API_KEY"]},
		}
		resp, err := http.PostForm("https://api.honeybadger.io/v1/deploys", params)
		if err == nil {
			fmt.Printf("Honeybadger: %v\n", resp)
		} else {
			fmt.Printf("Honeybadger Error: %v\n", err)
			fmt.Printf("Honeybadger Request: %v\n", params)
		}
	}
}
