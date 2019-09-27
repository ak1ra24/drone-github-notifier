package githubapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ak1ra24/drone-github-notifier/ci"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type Github struct {
	Client *github.Client
	Token  string
	Owner  string
	Repo   string
	PR     ci.PullRequest
}

func NewClient(owner, repo, gtoken string, pr ci.PullRequest) *Github {
	gtoken = strings.TrimLeft(gtoken, "$")
	token := os.Getenv(gtoken)
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	config := &Github{Client: client, Owner: owner, Repo: repo, Token: token, PR: pr}

	return config
}

func (g *Github) GetIssue() {

	events, _, err := g.Client.Issues.ListRepositoryEvents(context.Background(), g.Owner, g.Repo, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(events))
	for _, event := range events {
		fmt.Println(event.GetIssue())
	}
}

func (g *Github) GetPR() {

	prs, _, err := g.Client.PullRequests.List(context.Background(), g.Owner, g.Repo, nil)
	if err != nil {
		panic(err)
	}

	for _, pr := range prs {
		pr_comments, _, err := g.Client.PullRequests.ListComments(context.Background(), g.Owner, g.Repo, *pr.Number, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(pr_comments)
	}
}

func (g *Github) CreateIssue(title, body string, labels []string) error {
	issreq := &github.IssueRequest{Title: &title, Body: &body, Labels: &labels}
	_, _, err := g.Client.Issues.Create(context.Background(), g.Owner, g.Repo, issreq)
	if err != nil {
		return err
	}
	return nil
}

func (g *Github) PRComment(body string) error {

	fmt.Println("CHECK: ", g.Owner, g.Repo, g.PR.Number, g.PR.Reversion)
	if g.PR.Number != 0 {
		fmt.Println("In PR NUMBER")
		_, _, err := g.Client.Issues.CreateComment(context.Background(), g.Owner, g.Repo, g.PR.Number, &github.IssueComment{Body: &body})
		_, _, err = g.Client.Repositories.CreateComment(context.Background(), g.Owner, g.Repo, g.PR.Reversion, &github.RepositoryComment{Body: &body})
		if err != nil {
			return err
		}
	}
	if g.PR.Reversion != "" {
		fmt.Println("In Revision")
		prs, err := g.GetPRs()
		if err != nil {
			return err
		}
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println(prs)
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
		for _, pr := range prs {
			fmt.Println(pr)
			fmt.Println(&pr.Number, &pr.Title)
		}

		commits, err := g.List(g.PR.Reversion)
		if err != nil {
			return err
		}
		lastRevision, _ := g.lastOne(commits, g.PR.Reversion)
		g.PR.Reversion = lastRevision
		_, _, err = g.Client.Repositories.CreateComment(context.Background(), g.Owner, g.Repo, g.PR.Reversion, &github.RepositoryComment{Body: &body})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Github) GetPRs() ([]*github.PullRequest, error) {
	prs, _, err := g.Client.PullRequests.ListPullRequestsWithCommit(context.Background(), g.Owner, g.Repo, g.PR.Reversion, &github.PullRequestListOptions{})
	return prs, err
}

func (g *Github) List(revision string) ([]string, error) {
	if revision == "" {
		return []string{}, errors.New("no revision specified")
	}
	var s []string
	commits, _, err := g.Client.Repositories.ListCommits(context.Background(), g.Owner, g.Repo, &github.CommitsListOptions{SHA: revision})
	if err != nil {
		return s, err
	}
	for _, commit := range commits {
		s = append(s, *commit.SHA)
	}
	return s, nil
}

// Last returns the hash of the previous commit of the given commit
func (g *Github) lastOne(commits []string, revision string) (string, error) {
	if revision == "" {
		return "", errors.New("no revision specified")
	}
	if len(commits) == 0 {
		return "", errors.New("no commits")
	}
	// e.g.
	// a0ce5bf 2018/04/05 20:50:01 (HEAD -> master, origin/master)
	// 5166cfc 2018/04/05 20:40:12
	// 74c4d6e 2018/04/05 20:34:31
	// 9260c54 2018/04/05 20:16:20
	return commits[1], nil
}
