package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/semaphore"
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

	query := `org:github`
	resp, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var errs []error
	var mu sync.Mutex

	maxWorkers := 20
	sem := semaphore.NewWeighted(int64(maxWorkers))

	fmt.Printf("cloning %d repositories\n", len(resp.Repositories))
	for _, repo := range resp.Repositories {
		wg.Add(1)

		sem.Acquire(ctx, 1)

		repo := repo

		go func() {
			if err := clone(*repo.CloneURL); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}

			wg.Done()
			sem.Release(1)
		}()
	}

	wg.Wait()

	if len(errs) != 0 {
		fmt.Printf("Errors = %+v\n", errs)
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
