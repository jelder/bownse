# bownse
How Heroku's deploy webhooks were supposed to work: Ping a bunch of relevant services after every deploy. Currently supports Slack, Honeybadger, and NewRelic.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

### GitHub Deployments

Heroku doesn't expose (as far as I can tell) the mapping of app name to GitHub Repo and branch. Instead we have to configure this explicitly. Note that in this example, my app is _really_ called `myapp-production`, but due to environment variable identifier restrictions, the underscores have been replaced with hyphens.

```bash
heroku config:set GITHUB_REPO_myapp_production=jelder/myapp             --app my-bownse-instance
```

If your app deploying from something other than `master`, you can configure that, too.

```bash
heroku config:set GITHUB_BRANCH_MYAPP-PRODUCTION=production             --app my-bownse-instance
```

In the future, we may support the GitHub Deployment Status API. Here's how that will work.

```bash
heroku config:set GITHUB_USER=jelder                                    --app my-bownse-instance
heroku config:set GITHUB_TOKEN=3fade40de99c370e46551ee53d18d50184824ec4 --app my-bownse-instance
```

