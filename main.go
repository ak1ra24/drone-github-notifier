package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ak1ra24/drone-github-notifier/ci"
	"github.com/ak1ra24/drone-github-notifier/githubapi"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

type NotifierService struct {
	Ci       string   `yaml:"ci"`
	Notifier Notifier `yaml:"notifier"`
}

type Notifier struct {
	Github Github `yaml:"github"`
}

type Github struct {
	Token      string `yaml:"token"`
	Repository struct {
		Owner string `yaml:"owner"`
		Repo  string `yaml:"name"`
	} `yaml:"repository"`
}

func ReadYaml(filename string) NotifierService {
	// yamlを読み込む
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// structにUnmasrshal
	var notifier NotifierService
	err = yaml.Unmarshal(buf, &notifier)
	if err != nil {
		panic(err)
	}
	return notifier
}

func main() {
	args := os.Args
	if len(args) != 2 {
		os.Exit(1)
	}

	notifier := ReadYaml(args[1])
	ciname := notifier.Ci
	github_settings := notifier.Notifier.Github

	var ciservice ci.CiService
	var err error
	switch ciname {
	case "drone":
		ciservice, err = ci.Drone()
		if err != nil {
			panic(err)
		}
	case "":
		fmt.Errorf("Set CI Service")
	default:
		fmt.Errorf("Not Support")
	}

	if ciservice.Event == "pull_request" {
		pr := ciservice.PR
		fmt.Println(pr)

		client := githubapi.NewClient(github_settings.Repository.Owner, github_settings.Repository.Repo, github_settings.Token, pr)
		if terminal.IsTerminal(0) {
			fmt.Println("パイプ無し")
		} else {
			b, _ := ioutil.ReadAll(os.Stdin)
			fmt.Println("パイプで渡された内容:", string(b))
			if err := client.PRComment(string(b)); err != nil {
				panic(err)
			}
		}
	}
}
