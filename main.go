package main

import (
	"fmt"
	"mj/go-github/github"
	"os"
)

const rootApiUrl = "https://api.github.com"
const apiTokenName = "GITHUB_READ_API_TOKEN"

func main() {
	fmt.Println("Calling Github...")
	token, exists := os.LookupEnv(apiTokenName)
	if !exists {
		fmt.Printf("Error: %s not found", apiTokenName)
		return
	}
	client := github.NewGithubClient(token, rootApiUrl)
	prs, err := client.FetchContributions("mcMickJuice", "2024-08-01")

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

  fmt.Printf("PRs: %v\n", prs)

	fmt.Println("Program has completed")
}
