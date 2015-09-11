package main

import (
	. "github.com/jelder/bownse/env"
	"log"
	"regexp"
)

type Config struct {
	Secret        string
	ListenAddress string
	TargetURLs    []string
}

var (
	targetUrlRegexp = regexp.MustCompile(`^(\w+)_URL$`)
	config          Config
	ENV             EnvMap
)

func init() {
	ENV = MustLoadEnv()
	config.Secret = getSecret()
	config.ListenAddress = getListenAddress()
}

func getSecret() (secret string) {
	secret = ENV["SECRET"]
	if len(secret) < 30 {
		log.Fatal("SECRET is much too short; refusing to use it.")
	}
	return secret
}

func getListenAddress() (port string) {
	port = ENV["PORT"]
	if port == "" {
		return ":8080"
	} else {
		return ":" + port
	}
}
