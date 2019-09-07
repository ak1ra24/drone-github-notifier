package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ak1ra24/drone-github-notifier/githubapi"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	owner = flag.String("owner", "", "github owner")
	repo  = flag.String("repo", "", "github repo")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s [OPTIONS] ARGS...
Options\n`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	config := githubapi.NewClient(*owner, *repo)
	fmt.Println(terminal.IsTerminal(0))
	if terminal.IsTerminal(0) {
		fmt.Println("パイプ無し(FD値0)")
	} else {
		b, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println("パイプで渡された内容(FD値0以外):", string(b))
		drone_pr_num := os.Getenv("DRONE_PULL_REQUEST")
		if len(drone_pr_num) != 0 {
			config.PRComment(1, string(b))
		} else {
			fmt.Errorf("Not Setting DRONE_PULL_REQUEST")
		}
	}
}
