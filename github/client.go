package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SearchResponse[I any] struct {
	TotalCount int `json:"total_count"`
	Items      []I `json:"items"`
}

type SearchRepoResponseItem struct {
	Id     int    `json:"id"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
}

type SearchPullRequestResponseItem struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
}

type GithubClient struct {
	token, baseUrl string
}

func NewGithubClient(token, baseUrl string) GithubClient {
	return GithubClient{token, baseUrl}
}

func (c GithubClient) fetch(path string, method string, data interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
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
	path := "/search/repositories?q=org:shipt+segway+in:name&per_page=30"
	repoResponse := &SearchResponse[SearchRepoResponseItem]{}

	if err := c.fetch(path, http.MethodGet, repoResponse); err != nil {
		return nil, err
	}
	fmt.Printf("Total Count: %d\n", repoResponse.TotalCount)
	var repoNames []string
	for _, repo := range repoResponse.Items {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}

func (c GithubClient) FetchContributions(user, sinceDate string) ([]SearchPullRequestResponseItem, error) {
	path := fmt.Sprintf("/search/issues?q=is:pr+repo:shipt/segway-next+author:%s+created:>%s", user, sinceDate)
	searchResponse := &SearchResponse[SearchPullRequestResponseItem]{}
	if err := c.fetch(path, http.MethodGet, searchResponse); err != nil {
		return nil, err
	}
	fmt.Printf("Total Count: %d\n", searchResponse.TotalCount)
	return searchResponse.Items, nil
}
