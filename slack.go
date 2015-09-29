// https://api.slack.com/incoming-webhooks
// curl -X POST --data-urlencode 'payload={"text": "This is posted to <#general> and comes from *monkey-bot*.", "channel": "#general", "username": "monkey-bot", "icon_emoji": ":monkey_face:"}' https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/jelder/goenvmap"
	"net/http"
	"strconv"
)

func init() {
	if !SlackIsConfigured() {
		fmt.Println("Slack is not configured.")
	}
}

func SlackIsConfigured() bool {
	return ENV["SLACK_URL"] != ""
}

type slackMessage struct {
	Text        string            `json:"text,omitempty"`
	UserName    string            `json:"username,omitempty"`
	IconUrl     string            `json:"icon_url,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Attachments []slackAttachment `json:"attachments,omitempty"`
}

type slackAttachment struct {
	Fallback   string       `json:"fallback,omitempty"`
	Color      string       `json:"color,omitempty"`
	Pretext    string       `json:"pretext,omitempty"`
	AuthorName string       `json:"author_name,omitempty"`
	AuthorLink string       `json:"author_link,omitempty"`
	AuthorIcon string       `json:"author_icon,omitempty"`
	Title      string       `json:"title,omitempty"`
	TitleLink  string       `json:"title_link,omitempty"`
	Text       string       `json:"text,omitempty"`
	Fields     []slackField `json:"fields,omitempty"`
	ImageThumb string       `json:"image_thumb,omitempty"`
	ImageUrl   string       `json:"image_url,omitempty"`
	MarkdownIn []string     `json:"mrkdwn_in,ommitempty`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func SlackRequest(state *HerokuAppState) *http.Request {
	message := slackMessage{
		UserName:  "Heroku Deployment",
		IconUrl:   "https://d1ic07fwm32hlr.cloudfront.net/images/favicon.ico",
		IconEmoji: ":heroku:",
		Attachments: []slackAttachment{
			{
				Fallback:   fmt.Sprintf("%s deployed %s %s", state.User, state.App, state.Release),
				Color:      "#430098",
				AuthorName: state.User,
				Text:       fmt.Sprintf(" %s\n", state.GitLog),
				Title:      fmt.Sprintf("%s %s", state.App, state.Release),
				TitleLink:  state.URL(),
				Fields: []slackField{
					{
						Title: "Current Commit",
						Value: state.Head,
						Short: true,
					},
				},
			},
		},
	}

	if state.PrevHead != "" {
		field := slackField{
			Title: "Previous Commit",
			Value: state.PrevHead,
			Short: true,
		}
		message.Attachments[0].Fields = append(message.Attachments[0].Fields, field)
	}

	jsonStr, _ := json.MarshalIndent(message, "", "  ")
	fmt.Printf("%s\n", jsonStr)
	req, _ := http.NewRequest("POST", ENV["SLACK_URL"], bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonStr)))

	return req
}
