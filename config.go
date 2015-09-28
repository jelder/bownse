package main

import (
	. "github.com/jelder/goenvmap"
	"log"
)

type Config struct {
	Secret          string
	ListenAddress   string
	TargetURLs      []string
	HerokuAuthToken string
}

var (
	config Config
)

func init() {
	ENV = MustLoadEnv()
	config.Secret = getSecret()
	config.ListenAddress = getListenAddress()
	config.HerokuAuthToken = getHerokuAuthToken()
}

func getHerokuAuthToken() (authToken string) {
	authToken = ENV["HEROKU_AUTH_TOKEN"]
	if authToken == "" {
		log.Fatal("HEROKU_AUTH_TOKEN is not set; cannot fetch your Heroku apps' ENV vars")
	}
	return authToken
}

func getSecret() (secret string) {
	secret = ENV["SECRET_KEY"]
	if len(secret) < 30 {
		log.Fatal("SECRET_KEY is much too short; refusing to use it.")
	}
	return secret
}

func getListenAddress() string {
	return ":" + ENV.Get("PORT", "8080")
}
