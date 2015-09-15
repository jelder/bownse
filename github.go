// https://developer.github.com/v3/repos/deployments/

package main

import (
	"fmt"
	. "github.com/jelder/goenv"
)

func init() {
	if !GitHubIsConfigured() {
		fmt.Println("GitHub is not configured.")
	}
}

func GitHubIsConfigured() bool {
	return ENV["GITHUB_USER"] != "" && ENV["GITHUB_TOKEN"] != ""
}

func GitHubRepo(app string) string {
	return ENV[fmt.Sprint("GITHUB_REPO_", app)]
}
