package main

import (
	. "github.com/jelder/goenv"
	"log"
)

type Config struct {
	Secret        string
	ListenAddress string
	TargetURLs    []string
}

var (
	config Config
)

func init() {
	ENV = MustLoadEnv()
	config.Secret = getSecret()
	config.ListenAddress = getListenAddress()
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
