package github

type PullRequestResult struct {
	Title string
	Date  string
}

type PullRequestReviewOverview struct {
	PullRequests []PullRequestResult
	Reviews      []PullRequestResult
}

func (o PullRequestReviewOverview) PrCount() int {
	return len(o.PullRequests)
}

func (o PullRequestReviewOverview) ReviewCount() int {
	return len(o.Reviews)
}

type searchResponse[I any] struct {
	TotalCount int `json:"total_count"`
	Items      []I `json:"items"`
}

type searchRepoResponseItem struct {
	Id     int    `json:"id"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
}

type searchPullRequestResponseItem struct {
	Title     string `json:"title"`
	Number    int    `json:"number"`
	CreatedAt string `json:"created_at"`
}
