package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	user := flag.String("u", "", "Github username")
	key := flag.String("key", "", "Github access key")
	dir := flag.String("o", "", "Backup directory")
	forks := flag.Bool("forks", false, "Enable backup of forked repos")
	flag.Parse()

	if len(*user) == 0 {
		fmt.Println("Github username not specified")
		os.Exit(1)
	}

	if len(*dir) == 0 {
		fmt.Println("Backup directory not specified")
		os.Exit(1)
	}

	getTokenSource := func() oauth2.TokenSource {
		if len(*key) > 0 {
			return oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: *key},
			)
		}
		return nil
	}

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, getTokenSource())
	client := gh.NewClient(tc)

	userRepos, _, err := client.Repositories.List(ctx, *user, nil)
	if err != nil {
		panic(err)
	}

	for _, repo := range userRepos {
		if *repo.Fork && !*forks {
			continue
		}

		fmt.Println(*repo.Name)
	}
}
