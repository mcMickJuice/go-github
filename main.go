package main

import (
	"fmt"
	"mj/go-github/github"
)

const rootApiUrl = "https://api.github.com"
func main() {
	fmt.Println("Hello!")
	repoNames, err := github.FetchRepos(rootApiUrl)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	for _, name := range repoNames {
		fmt.Println(name)
	}

	fmt.Println("Program has completed")
}
