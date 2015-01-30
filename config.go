package main

import (
	"github.com/jelder/bownse/env"
	"log"
	neturl "net/url"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	Secret        string
	ListenAddress string
	TargetURLs    []string
}

var (
	targetUrlRegexp = regexp.MustCompile(`^(\w+)_URL$`)
	config          Config
)

func init() {
	ENV := MustLoadEnv()
	config.Secret = getSecret()
	config.ListenAddress = ENV.Get
	config.TargetURLs = getTargetURLs
}

func getSecret() (secret string) {
	secret = ENV["SECRET"]
	if len(secret) < 30 {
		log.Fatal("SECRET is much too short; refusing to use it.")
	}
	return secret
}

func getListenAddress() string {
	string = ENV["PORT"]
	if string == "" {
		return ":8080"
	} else {
		return ":" + string
	}
}

func getTargetURLs() (urls []string) {
	for _, e := range ENV {
		pair := strings.Split(e, "=")
		if url := CheckEnvUrl(pair); url != "" {
			urls = append(urls, pair[1])
		}
	}
	return urls
}

func CheckEnvUrl(pair []string) string {
	if targetUrlRegexp.MatchString(pair[0]) {
		url, err := neturl.Parse(pair[1])
		if err != nil {
			return ""
		}
		if url.IsAbs() && (url.Scheme == "http" || url.Scheme == "https") {
			return url.String()
		}
	}
	return ""
}
