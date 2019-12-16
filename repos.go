package main

import (
	"context"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getRepos(user, key string, forks bool) ([]string, error) {
	getTokenSource := func() oauth2.TokenSource {
		if len(key) > 0 {
			return oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: key},
			)
		}
		return nil
	}

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, getTokenSource())
	client := gh.NewClient(tc)

	userRepos, _, err := client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, repo := range userRepos {
		if *repo.Fork && !forks {
			continue
		}

		results = append(results, user+"/"+*repo.Name)
	}

	return results, nil
}
