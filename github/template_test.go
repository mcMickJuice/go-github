package github

import (
	"bytes"
	"testing"
)

func TestTemplate(t *testing.T) {
	overview := PullRequestReviewOverview{
		PullRequests: []PullRequestResult{
			PullRequestResult{"PR 1", "date"},
		},
		Reviews: []PullRequestResult{
			PullRequestResult{"Review 1", "date"},
			PullRequestResult{"Review 2", "date"},
		},
	}
	buf := bytes.Buffer{}
	err := RenderGithubContribution(&buf, overview)
	want := "PRs: 1\nReviews: 2\n"

	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
