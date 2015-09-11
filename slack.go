// https://api.slack.com/incoming-webhooks
// curl -X POST --data-urlencode 'payload={"text": "This is posted to <#general> and comes from *monkey-bot*.", "channel": "#general", "username": "monkey-bot", "icon_emoji": ":monkey_face:"}' https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

package main

var (
	slack = make(chan *HerokuWebhookPayload)
)

func init() {
	go handleSlack()
}

func handleSlack() {
	for {
		payload := <-slack
		// TODO
		// resp, err := http.PostForm(ENV["SLACK_URL"], url.Values{}
	}
}
