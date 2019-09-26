package ci

import (
	"os"
	"strconv"
)

type CiService struct {
	PR  PullRequest
	URL string
}

type PullRequest struct {
	Reversion string
	Number    int
	Body      string
}

func Drone() (ci CiService, err error) {
	ci.PR.Number = 0
	ci.PR.Reversion = os.Getenv("DRONE_COMMIT_SHA")
	ci.URL = os.Getenv("DRONE_BUILD_LINK")
	pr := os.Getenv("DRONE_PULL_REQUEST")
	if pr == "" {
		return ci, err
	}
	ci.PR.Number, err = strconv.Atoi(pr)
	return ci, err
}
