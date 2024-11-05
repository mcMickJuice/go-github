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
