// https://docs.newrelic.com/docs/apm/new-relic-apm/maintenance/deployment-notifications#examples
// curl -H "x-api-key:REPLACE_WITH_YOUR_API_KEY" -d "deployment[app_name]=REPLACE_WITH_YOUR_APP_NAME" -d "deployment[description]=This is an app id deployment" https://api.newrelic.com/deployments.xml

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var (
	newrelic = make(chan *HerokuWebhookPayload)
)

func init() {
	if NewRelicIsConfigured() {
		go handleNewRelic()
	} else {
		fmt.Println("NewRelic is not fully configured.")
	}
}

func NewRelicIsConfigured() bool {
	return ENV["NEW_RELIC_APP_NAME"] != "" && ENV["NEW_RELIC_ID"] != "" && ENV["NEW_RELIC_API_KEY"] != ""

}

func handleNewRelic() {
	for {
		payload := <-newrelic
		apiUrl := "https://api.newrelic.com/"
		resource := "deployments.xml"
		params := url.Values{
			"deployment[app_name]":       {ENV["NEW_RELIC_APP_NAME"]},
			"deployment[application_id]": {ENV["NEW_RELIC_ID"]},
			"deployment[user]":           {payload.User},
			"deployment[description]":    {""},
			"deployment[changelog]":      {""},
			"deployment[revision]":       {payload.Head},
		}

		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = resource
		urlStr := fmt.Sprintf("%v", u)

		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(params.Encode()))
		r.Header.Add("x-api-key", ENV["NEW_RELIC_API_KEY"])
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

		resp, err := client.Do(r)
		if err == nil {
			fmt.Printf("NewRelic: %v\n", resp)
		} else {
			fmt.Printf("NewRelic Error: %v\n", err)
			fmt.Printf("NewRelic Params: %v\n", params)
		}

	}
}
