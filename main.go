package main

import (
	"flag"
	"fmt"
	"mj/go-github/github"
	"os"
)

const rootApiUrl = "https://api.github.com"
const apiTokenName = "GITHUB_READ_API_TOKEN"

func main() {
	token, exists := os.LookupEnv(apiTokenName)

	if !exists {
		fmt.Printf("Error: %s not found", apiTokenName)
		return
	}
	client := github.NewGithubClient(token, rootApiUrl)
	user := flag.String("user", "", "Github User Handle")
	from := flag.String("from", "2024-08-01", "Contributions since Date YYYY-MM-DD. default: 2024-08-01")
  flag.Parse()

	if *user == "" {
		fmt.Println("user not provided")
		return
	}
  fmt.Printf("Fetching Github Info for %s from %s\n", *user, *from)
	prOverview, err := client.FetchContributions(*user, *from)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf("PR Count: %d\nReview Count: %d\n", len(prOverview.PullRequests), len(prOverview.Reviews))

	fmt.Println("Program has completed")
}
