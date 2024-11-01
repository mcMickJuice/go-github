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
	Items      []SearchRepoResponseItem
}

type SearchRepoResponseItem struct {
	Id     int    `json:"id"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
}

type GithubClient struct {
	token, baseUrl string
}

func NewGithubClient(token, baseUrl string) GithubClient {
	return GithubClient{token, baseUrl}
}

func (c GithubClient) fetch(url string, method string, data interface{}) error {

	request, _ := http.NewRequest(method, url, nil)
	client := &http.Client{}
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}
	return nil
}

// fetch all repos available to user
func (c GithubClient) FetchRepos() ([]string, error) {
  url := fmt.Sprintf("%s/search/repositories?q=org:shipt+segway+in:name&per_page=30", c.baseUrl)
	repoResponse := &SearchRepoResponse{}

	if err := c.fetch(url, http.MethodGet, repoResponse); err != nil {
		return nil, err
	}
	fmt.Printf("Total Count: %d\n", repoResponse.TotalCount)
	var repoNames []string
	for _, repo := range repoResponse.Items {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}

func (c GithubClient) FetchContributions(user string) error {
	return nil
}
