package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	query := `org:github language:go`
	resp, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	fmt.Printf("cloning %d repositories\n", len(resp.Repositories))
	for _, repo := range resp.Repositories {
		wg.Add(1)

		repo := repo

		go func() {
			if err := clone(*repo.CloneURL); err != nil {
				panic(err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("finished cloning")
}

func clone(url string) error {
	fmt.Println("cloning", url)
	c := exec.Command("git", "clone", url)
	_, err := c.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
