# bownse
How Heroku's deploy webhooks were supposed to work: Ping a bunch of relevant services after every deploy. Currently supports Slack, Honeybadger, and NewRelic.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

### GitHub Deployments

Heroku doesn't expose (as far as I can tell) the mapping of app name to GitHub Repo and branch. Instead we have to configure this explicitly.

```bash
heroku config:set GITHUB_REPO_myapp-production=jelder/myapp             --app my-bownse-instance
heroku config:set GITHUB_BRANCH_myapp-production=master                 --app my-bownse-instance # Defaults to `master`
heroku config:set GITHUB_USER=jelder                                    --app my-bownse-instance
heroku config:set GITHUB_TOKEN=3fade40de99c370e46551ee53d18d50184824ec4 --app my-bownse-instance
```

