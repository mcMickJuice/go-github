package main

import (
	"fmt"
	"mj/go-github/github"
	"os"
)

const rootApiUrl = "https://api.github.com"
const apiTokenName = "GITHUB_READ_API_TOKEN"

func main() {
	fmt.Println("Hello!")
	token, exists := os.LookupEnv(apiTokenName)
	if !exists {
		fmt.Printf("Error: %s not found", apiTokenName)
		return
	}
	repoNames, err := github.FetchRepos(rootApiUrl, token)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	for _, name := range repoNames {
		fmt.Println(name)
	}

	fmt.Println("Program has completed")
}
