package github

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed contributions.txt
var contributionTextTemplate string

func RenderGithubContribution(w io.Writer, overview PullRequestReviewOverview) error {
	contributionTemplate := template.New("contribution")
	contributionTemplate, err := contributionTemplate.Parse(contributionTextTemplate)

	if err != nil {
		return err
	}

	if err := contributionTemplate.Execute(w, overview); err != nil {
		return err
	}

	return nil
}
