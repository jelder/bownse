// https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook
// The parameters included in the request are the same as the variables available in the hook message: app, user, url, head, head_long, git_log and release. See below for their descriptions.
// This is an example payload:
// curl -X POST http://localhost:8080/$SECRET -d app=secure-woodland-9775 -d user=example%40example.com -d url=http://secure-woodland-9775.herokuapp.com -d head=4f20bdd -d head_long=4f20bdd -d prev_head= -d git_log=%20%20*%20Michael%20Friis%3A%20add%20bar -d release=v7

package main

import (
	"regexp"
)

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
}

var (
	stagingRegexp    = regexp.MustCompile(`staging`)
	productionRegexp = regexp.MustCompile(`production|\bprod\b`)
)

func (payload *HerokuWebhookPayload) Environment() string {
	switch {
	case stagingRegexp.MatchString(payload.App):
		return "staging"
	case productionRegexp.MatchString(payload.App):
		return "production"
	default:
		return "development"
	}
}
