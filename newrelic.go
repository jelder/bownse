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

func NewRelicRequest(state *HerokuAppState) *http.Request {
	urlStr := "https://api.newrelic.com/deployments.xml"
	params := url.Values{
		"deployment[app_name]":       {state.Env["NEW_RELIC_APP_NAME"]},
		"deployment[application_id]": {state.Env["NEW_RELIC_ID"]},
		"deployment[user]":           {state.User},
		"deployment[description]":    {fmt.Sprintf("%s %s", state.App, state.Release)},
		"deployment[changelog]":      {fmt.Sprintf("  %s", state.GitLog)},
		"deployment[revision]":       {state.Head},
	}

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(params.Encode()))
	req.Header.Add("x-api-key", state.Env["NEW_RELIC_LICENSE_KEY"])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	return req
}
