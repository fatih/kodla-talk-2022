package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/google/go-github/v44/github"
)

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	query := `org:github language:go`
	resp, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		panic(err)
	}

	for _, repo := range resp.Repositories {
		go func() {
			if err := clone(*repo.CloneURL); err != nil {
				panic(err)
			}
		}()
	}

	fmt.Println("finished cloning")
}

func clone(url string) error {
	fmt.Println("cloning", url)
	c := exec.Command("git", "clone", url)
	out, err := c.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Printf("string(out = %+v\n", string(out))
	return nil
}
