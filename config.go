package main

import (
	"log"
	neturl "net/url"
	"os"
	"regexp"
	"strings"
)

var (
	targetUrlRegexp = regexp.MustCompile(`^(\w+)_URL$`)
	Secret          string
	ListenAddress   string
	TargetUrls      []string
)

func init() {
	Secret = GetSecret()
	ListenAddress = GetListenAddress()
	TargetUrls = GetTargetUrls()
}

func GetSecret() (secret string) {
	secret = os.Getenv("SECRET")
	if len(secret) < 30 {
		log.Fatal("SECRET is much too short; refusing to use it.")
	}
	return secret
}

func GetListenAddress() string {
	string := os.Getenv("PORT")
	if string == "" {
		return ":8080"
	} else {
		return ":" + string
	}
}

func GetTargetUrls() (urls []string) {
	for _, e := range os.Environ() {
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
