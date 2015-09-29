# bownse
[![Circle CI](https://circleci.com/gh/jelder/bownse.svg?style=svg)](https://circleci.com/gh/jelder/bownse)

How Heroku's [deploy webhooks](https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook) were supposed to work: Ping a bunch of relevant services after every deploy. Bownse sits between Heroku and your addons, sending artisnal webhooks to your team chat, application performance monitors, and error trackers.

Currently supports Slack, Honeybadger, and NewRelic.

You can launch your own private Bownse instance from this button.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## Setup

Once you've deployed your own copy of Bownse, configure all of your Heroku apps to use it. The URL will be the URL of your Bownse instance followed by the value `SECRET_KEY` ENV var. For example, it might look a lot like this:

```
https://boundless-bownse.herokuapp.com/074585ce6ef2d8e457d31fc4af098bbdae039c640f041184f9b2488d60e19012
```

Heroku doesn't expose the upstream GitHub respository name, so you will have to configure it manually as an ENV var for each of your apps.

```bash
heroku config:set GITHUB_REPO=myname/myapp --app myapp
```

Finally, you must tell your Bownse instance how to authenticate against Heroku to get the information it needs. It exclusively uses the `config-vars` endpoint and never modifies anything. More details here: https://devcenter.heroku.com/articles/platform-api-reference#config-vars

```bash
heroku config:set HEROKU_AUTH_TOKEN=$(heroku auth:token) --app my-bownse-instance
```

## Slack

Bownse will figure everything out from your app's ENV vars, with one exception: Slack. Create an incoming webhook at https://boundless.slack.com/services/new/incoming-webhook, and then tell your Bownse instance about it by setting the `SLACK_URL` ENV var.

If you don't use Slack, I envy you. It's a pretty mediocre system especially for developers. 

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
