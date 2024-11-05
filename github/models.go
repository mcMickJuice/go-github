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
