// http://docs.honeybadger.io/article/38-deployment-tracking-on-heroku
// https://api.honeybadger.io/v1/deploys?deploy[environment]=production&deploy[local_username]={{user}}&deploy[revision]={{head}}&api_key=asdf

package main

import (
	"bytes"
	"fmt"
	. "github.com/jelder/goenvmap"
	"net/http"
	"net/url"
	"strconv"
)

func HoneybadgerRequest(payload *HerokuWebhookPayload) *http.Request {
	urlStr := "https://api.honeybadger.io/v1/deploys"
	params := url.Values{
		"deploy[environment]":    {payload.Environment()},
		"deploy[local_username]": {payload.User},
		"deploy[revision]":       {payload.Head},
		"api_key":                {payload.Env["HONEYBADGER_API_KEY"]},
	}

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	return req
}
