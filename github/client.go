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

type SearchRepoResponse struct {
TotalCount int `json:"total_count"`
  Items []SearchRepoResponseItem
}

type SearchRepoResponseItem struct {
	Id     int    `json:"id"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
}

// fetch all repos available to user
func FetchRepos(baseUrl, token string) ([]string, error) {
  url := fmt.Sprintf("%s/search/repositories", baseUrl)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	client := &http.Client{}
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	repoResponse := SearchRepoResponse{}
	json.Unmarshal(body, &repoResponse)
	if err != nil {
		return nil, err
	}
	var repoNames []string
	for _, repo := range repoResponse.Items{
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}
