package github

import (
	"fmt"
	"strings"
)

type GithubSearchQuery struct {
	terms []Query
}

func (q *GithubSearchQuery) Build() string {
	var terms []string
	for _, term := range q.terms {
		terms = append(terms, term.Build())
	}

	j := fmt.Sprintf("q=%s", strings.Join(terms, " "))
	return strings.ReplaceAll(j, " ", "+")
}

func (q *GithubSearchQuery) Add(term Query) *GithubSearchQuery {
	q.terms = append(q.terms, term)
	return q
}

type Query interface {
	Build() string
}

// org:{orgname}
type OrgQuery struct {
	value string
}

func (q OrgQuery) Build() string {
	return fmt.Sprintf("org:%s", q.value)
}

// {repoName} in:name
type RepoNameQuery struct {
	repoName string
}

func (q RepoNameQuery) Build() string {
	return fmt.Sprintf("%s in:name", q.repoName)
}

//is:pr
type IsPrQuery struct{}

func (q IsPrQuery) Build() string {
	return "is:pr"
}

// repo:{org/repo}
type RepoIssueQuery struct {
	org, repo string
}

func (q RepoIssueQuery) Build() string {
	return fmt.Sprintf("repo:%s/%s", q.org, q.repo)
}

type PrInteractionQuery struct {
	isAuthor, negation bool
	userName           string
}

func (q PrInteractionQuery) Build() string {
	queryPart := ""
	if q.isAuthor {
		queryPart = "author"
	} else {
		queryPart = "reviewed-by"
	}

	if q.negation {
		queryPart = "-" + queryPart
	}

	return fmt.Sprintf("%s:%s", queryPart, q.userName)
}

type CreatedAfterQuery struct {
	date string
}

func (q CreatedAfterQuery) Build() string {
	return fmt.Sprintf("created:>%s", q.date)
}
