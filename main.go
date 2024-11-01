package main

import (
	"fmt"
	"mj/go-github/github"
	"os"
)

const rootApiUrl = "https://api.github.com"
const apiTokenName = "GITHUB_READ_API_TOKEN"

// TODO - parallelize search calls, accept contributor and since Date, github query builder
func main() {
	fmt.Println("Calling Github...")
	token, exists := os.LookupEnv(apiTokenName)
	if !exists {
		fmt.Printf("Error: %s not found", apiTokenName)
		return
	}
	client := github.NewGithubClient(token, rootApiUrl)
	prOverview, err := client.FetchContributions("dcheng1290", "2024-08-01")

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf("PR Count: %d\nReview Count: %d\n", len(prOverview.PullRequests), len(prOverview.Reviews))

	fmt.Println("Program has completed")
}
