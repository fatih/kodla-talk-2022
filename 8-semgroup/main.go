package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/semgroup"
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
		ListOptions: github.ListOptions{PerPage: 100},
	}

	query := `org:github language:go`
	resp, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		panic(err)
	}

	const maxWorkers = 20
	s := semgroup.NewGroup(context.Background(), maxWorkers)

	fmt.Printf("cloning %d repositories\n", len(resp.Repositories))
	for _, repo := range resp.Repositories {
		repo := repo

		s.Go(func() error { return clone(*repo.CloneURL) })
	}

	if err := s.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("finished cloning")
}

func clone(url string) error {
	fmt.Println("cloning", url)
	c := exec.Command("git", "clone", "--depth=1", url)
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cloning failed: %s, out: %s", err, string(out))
	}

	fmt.Println("finished", url)
	return nil
}
