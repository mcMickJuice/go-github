package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Title     string `json:"title"`
	Number    int    `json:"number"`
	CreatedAt string `json:"created_at"`
}

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

type GithubSearchQuery struct {
	terms []Query
}

func (q *GithubSearchQuery) Build() string {
  var terms []string
	for _, term := range q.terms {
		terms = append(terms, term.Build())
	}

  j := fmt.Sprintf("q=%s", strings.Join(terms, " " ))
	return strings.ReplaceAll(j, " ", "+")
}

func (q *GithubSearchQuery) Add(term Query) *GithubSearchQuery {
	q.terms = append(q.terms, term)
	return q
}

type Query interface {
	Build() string
}

type OrgQuery struct {
	value string
}

func (q OrgQuery) Build() string {
	return fmt.Sprintf("org:%s", q.value)
}

type RepoNameQuery struct {
	repoName string
}

func (q RepoNameQuery) Build() string {
	return fmt.Sprintf("%s in:name", q.repoName)
}

// fetch all repos available to user
func (c GithubClient) FetchRepos() ([]string, error) {
	query := GithubSearchQuery{}
	q := query.Add(OrgQuery{"shipt"}).Add(RepoNameQuery{"segway"}).Build()
	path := fmt.Sprintf("/search/repositories?%s&per_page=30", q)
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

type PullRequestResult struct {
	Title string
	Date  string
}

type PullRequestReviewOverview struct {
	PullRequests []PullRequestResult
	Reviews      []PullRequestResult
}

func (c GithubClient) FetchContributions(user, sinceDate string) (PullRequestReviewOverview, error) {
	prPath := fmt.Sprintf("/search/issues?per_page=100&q=is:pr+repo:shipt/segway-next+author:%s+created:>%s", user, sinceDate)
	prSearchResponse := &SearchResponse[SearchPullRequestResponseItem]{}

	// this looks to be pulled in PRs where I leave comments...
	reviewPath := fmt.Sprintf("/search/issues?per_page=100&q=is:pr+repo:shipt/segway-next+reviewed-by:%s+-author:%s+created:>%s", user, user, sinceDate)
	reviewSearchResponse := &SearchResponse[SearchPullRequestResponseItem]{}

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
