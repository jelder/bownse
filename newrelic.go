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

func init() {
	if !NewRelicIsConfigured() {
		fmt.Println("NewRelic is not fully configured.")
	}
}

func NewRelicIsConfigured() bool {
	return ENV["NEW_RELIC_APP_NAME"] != "" && ENV["NEW_RELIC_ID"] != "" && ENV["NEW_RELIC_API_KEY"] != ""

}

func NewRelicRequest(payload *HerokuWebhookPayload) *http.Request {
	urlStr := "https://api.newrelic.com/deployments.xml"
	params := url.Values{
		"deployment[app_name]":       {ENV["NEW_RELIC_APP_NAME"]},
		"deployment[application_id]": {ENV["NEW_RELIC_ID"]},
		"deployment[user]":           {payload.User},
		"deployment[description]":    {fmt.Sprintf("%s %s", payload.App, payload.Release)},
		"deployment[changelog]":      {payload.GitLog},
		"deployment[revision]":       {payload.Head},
	}

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(params.Encode()))
	req.Header.Add("x-api-key", ENV["NEW_RELIC_API_KEY"])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	return req
}
