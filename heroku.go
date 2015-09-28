// https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook
// The parameters included in the request are the same as the variables available in the hook message: app, user, url, head, head_long, git_log and release. See below for their descriptions.
// This is an example payload:
// curl -X POST http://localhost:8080/$SECRET -d app=secure-woodland-9775 -d user=example%40example.com -d url=http://secure-woodland-9775.herokuapp.com -d head=4f20bdd -d head_long=4f20bdd -d prev_head= -d git_log=%20%20*%20Michael%20Friis%3A%20add%20bar -d release=v7

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"regexp"
)

var (
	decoder          = schema.NewDecoder()
	stagingRegexp    = regexp.MustCompile(`staging`)
	productionRegexp = regexp.MustCompile(`production|\bprod\b`)
	heads            = make(map[string]string)
	client           = &http.Client{}
)

func init() {
	decoder.IgnoreUnknownKeys(true)
}

type HerokuAppEnv map[string]string

type HerokuWebhookPayload struct {
	App      string `schema:"app"`
	User     string `schema:"user"`
	Url      string `schema:"url"`
	Head     string `schema:"head"`
	HeadLong string `schema:"head_long"`
	PrevHead string `schema:"prev_head"`
	GitLog   string `schema:"git_log"`
	Release  string `schema:"release"`
	AppUUID  string `schema:"app_uuid"`
	Env      map[string]string
}

func ParseWebhook(r *http.Request) (payload *HerokuWebhookPayload, err error) {
	payload = new(HerokuWebhookPayload)
	err = decoder.Decode(payload, r.PostForm)
	if err != nil {
		fmt.Printf("Recieved Heroku Deploy Webhook: %+v\n", payload)
	}
	if payload.PrevHead == "" {
		payload.PrevHead = heads[payload.App]
		heads[payload.App] = payload.Head
	}
	payload.FetchEnv(config.HerokuAuthToken)
	return payload, err
}

func (payload *HerokuWebhookPayload) Environment() string {
	if payload.Env["RAILS_ENV"] != "" {
		return payload.Env["RAILS_ENV"] != ""
	}
	if payload.Env["RACK_ENV"] != "" {
		return payload.Env["RACK_ENV"] != ""
	}
	switch {
	case stagingRegexp.MatchString(payload.App):
		return "staging"
	case productionRegexp.MatchString(payload.App):
		return "production"
	default:
		return "development"
	}
}

// Return a GitHub compare URL if the repository is configured, otherwise just return the plain URL.
func (payload *HerokuWebhookPayload) URL() (url string) {
	repo := GitHubRepo(payload.App)
	if repo == "" {
		url = payload.Url
	} else {
		url = "https://github.com/" + repo
		if payload.PrevHead != "" {
			url = fmt.Sprint(url, "/compare/", payload.PrevHead, "...", payload.HeadLong)
		}
	}
	return url
}

func (payload *HerokuWebhookPayload) FetchEnv(string authTok) {
	payload.Env = MustGetHerokuAppEnv(payload.App, authToken)
}

func MustGetHerokuAppEnv(appName string, authToken string) (appEnv *HerokuAppEnv) {
	req, err := http.NewRequest("GET", "https://api.heroku.com/apps/"+appName+"/config-vars", nil)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	json.NewDecoder(resp.Body).Decode(&appEnv)
	return appEnv
}
