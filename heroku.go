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
)

func init() {
	decoder.IgnoreUnknownKeys(true)
}

type HerokuAppEnv map[string]string

type HerokuAppState struct {
	App      string `schema:"app"`
	User     string `schema:"user"`
	Url      string `schema:"url"`
	Head     string `schema:"head"`
	HeadLong string `schema:"head_long"`
	PrevHead string `schema:"prev_head"`
	GitLog   string `schema:"git_log"`
	Release  string `schema:"release"`
	AppUUID  string `schema:"app_uuid"`
	Env      HerokuAppEnv
}

func ParseWebhook(r *http.Request) (state *HerokuAppState, err error) {
	state = new(HerokuAppState)
	err = decoder.Decode(state, r.PostForm)
	if err != nil {
		fmt.Printf("Recieved Heroku Deploy Webhook: %+v\n", state)
	}
	if state.PrevHead == "" {
		state.PrevHead = heads[state.App]
		heads[state.App] = state.Head
	}
	state.FetchEnv(config.HerokuAuthToken)
	return state, err
}

func (state *HerokuAppState) Environment() string {
	if state.Env["RAILS_ENV"] != "" {
		return state.Env["RAILS_ENV"]
	} else if state.Env["RACK_ENV"] != "" {
		return state.Env["RACK_ENV"]
	} else {
		switch {
		case stagingRegexp.MatchString(state.App):
			return "staging"
		case productionRegexp.MatchString(state.App):
			return "production"
		default:
			return "development"
		}
	}
}

// Return a GitHub compare URL if the repository is configured, otherwise just return the plain URL.
func (state *HerokuAppState) URL() (url string) {
	if state.Env["GITHUB_REPO"] == "" {
		url = state.Url
	} else {
		url = "https://github.com/" + state.Env["GITHUB_REPO"]
		if state.PrevHead != "" {
			url = fmt.Sprint(url, "/compare/", state.PrevHead, "...", state.HeadLong)
		}
	}
	return url
}

func (state *HerokuAppState) FetchEnv(authToken string) {
	state.Env = MustGetHerokuAppEnv(state.App, authToken)
}

// https://devcenter.heroku.com/articles/platform-api-reference#config-vars
func MustGetHerokuAppEnv(appName string, authToken string) (appEnv HerokuAppEnv) {
	req, err := http.NewRequest("GET", "https://api.heroku.com/apps/"+appName+"/config-vars", nil)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("User-Agent", "Bownse: The Heroku Webhook Multiplexer")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	json.NewDecoder(resp.Body).Decode(&appEnv)
	return appEnv
}
