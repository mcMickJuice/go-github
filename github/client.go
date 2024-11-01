package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// curl -L \
//   -H "Accept: application/vnd.github+json" \
//   -H "Authorization: Bearer <YOUR-TOKEN>" \
//   -H "X-GitHub-Api-Version: 2022-11-28" \
//   https://api.github.com/orgs/ORG/repos

const rootApiUrl = "https://api.github.com"

type RepoResponse struct {
	Id     int    `json:"id"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
}

// fetch all repos available to user
func FetchRepos() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/orgs/shipt/repos", rootApiUrl))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	repoResponse := []RepoResponse{}
	json.Unmarshal(body, &repoResponse)
	if err != nil {
		return nil, err
	}
	var repoNames []string
	for _, repo := range repoResponse {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}
