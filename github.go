// https://developer.github.com/v3/repos/deployments/

package main

import (
	"fmt"
	. "github.com/jelder/goenv"
	"strings"
)

func GitHubRepo(app string) string {
	return ENV[fmt.Sprint("GITHUB_REPO_", strings.Replace(app, "-", "_", -1))]
}
