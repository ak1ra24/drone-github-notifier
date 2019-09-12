package githubapi

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type Github struct {
	Client *github.Client
	Owner  string
	Repo   string
}

func NewClient(owner, repo string) *Github {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	config := &Github{Client: client, Owner: owner, Repo: repo}

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

func (g *Github) PRComment(number int, body string) error {

	comment := &github.IssueComment{Body: &body}
	_, _, err := g.Client.Issues.CreateComment(context.Background(), g.Owner, g.Repo, number, comment)
	if err != nil {
		return err
	}

	return nil
}

func (g *Github) CreateIssue(title, body string, labels []string) error {
	issreq := &github.IssueRequest{Title: &title, Body: &body, Labels: &labels}
	_, _, err := g.Client.Issues.Create(context.Background(), g.Owner, g.Repo, issreq)
	if err != nil {
		return err
	}
	return nil
}
