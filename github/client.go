package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GithubClient struct {
	token, baseUrl string
}

func NewGithubClient(token, baseUrl string) GithubClient {
	return GithubClient{token, baseUrl}
}

func (c GithubClient) fetch(path string, method string, data interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	// when spaces in request, this is likely causing an error
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("fetch Error: StatusCode %d, url %s", resp.StatusCode, url))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	return nil
}

// fetch all repos available to user
func (c GithubClient) FetchRepos() ([]string, error) {
	query := GithubSearchQuery{queryParts: []Query{
		OrgQuery{"shipt"},
		RepoNameQuery{"segway"},
	}}
	path := fmt.Sprintf("/search/repositories?%s&per_page=30", query.Build())
	repoResponse := &searchResponse[searchRepoResponseItem]{}

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

func (c GithubClient) FetchContributions(user, sinceDate string) (PullRequestReviewOverview, error) {
	prQuery := GithubSearchQuery{queryParts: []Query{
		IsPrQuery{},
		RepoIssueQuery{org: "shipt", repo: "segway-next"},
		PrInteractionQuery{isAuthor: true, userName: user},
		CreatedAfterQuery{sinceDate},
	}}
	prPath := fmt.Sprintf("/search/issues?per_page=100&%s", prQuery.Build())
	prSearchResponse := &searchResponse[searchPullRequestResponseItem]{}

	// this looks to be pulled in PRs where I leave comments...
	reviewQuery := GithubSearchQuery{queryParts: []Query{
		IsPrQuery{},
		RepoIssueQuery{org: "shipt", repo: "segway-next"},
		PrInteractionQuery{isAuthor: false, userName: user},
		PrInteractionQuery{isAuthor: true, negation: true, userName: user},
		CreatedAfterQuery{sinceDate},
	}}
	reviewPath := fmt.Sprintf("/search/issues?per_page=100&%s", reviewQuery.Build())
	reviewSearchResponse := &searchResponse[searchPullRequestResponseItem]{}

	overview := PullRequestReviewOverview{}
	// parallelize this
	if err := c.fetch(prPath, http.MethodGet, prSearchResponse); err != nil {
		return overview, err
	}
	if err := c.fetch(reviewPath, http.MethodGet, reviewSearchResponse); err != nil {
		return overview, err
	}
	for _, item := range prSearchResponse.Items {
		overview.PullRequests = append(overview.PullRequests, PullRequestResult{item.Title, item.CreatedAt})
	}
	for _, item := range reviewSearchResponse.Items {
		overview.Reviews = append(overview.Reviews, PullRequestResult{item.Title, item.CreatedAt})
	}

	return overview, nil
}
