package github

import (
	"embed"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

//go:embed testdata/*
var testDataFs embed.FS

func TestClient(t *testing.T) {
	repoPath := "/search/repositories"
	t.Run("FetchRepos returns list of repo names", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != repoPath {
				t.Errorf("Expected path %s, Got path %s", repoPath, r.URL.Path)
			}

			wantQuery := url.Values{
				"q":        {"org:shipt segway in:name"},
				"per_page": {"30"},
			}
			if gotQuery := r.URL.Query(); !reflect.DeepEqual(gotQuery, wantQuery) {
				t.Errorf("Expected query %s, Got query %s", wantQuery, gotQuery)
			}

			if authHeader := r.Header.Get("Authorization"); authHeader != "Bearer token" {
				t.Errorf("Expected Header: %q, got %q", "Bearer token", authHeader)
			}
			searchRepoJson, _ := testDataFs.ReadFile("testdata/search_repo.json")
			w.WriteHeader(http.StatusOK)
			w.Write(searchRepoJson)
		}))
		defer server.Close()

		client := NewGithubClient("token", server.URL)
		want := []string{"hey", "hi"}
		got, err := client.FetchRepos()

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("Fetch Contributions", func(t *testing.T) {
		searchPrUrl := "per_page=100&q=is:pr+repo:shipt/segway-next+author:mickjuice+created:>2024-08-01"
		searchPrQuery, _ := url.ParseQuery(searchPrUrl)
		searchReviewsUrl := "per_page=100&q=is:pr+repo:shipt/segway-next+reviewed-by:mickjuice+-author:mickjuice+created:>2024-08-01"
		searchReviewsQuery, _ := url.ParseQuery(searchReviewsUrl)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if reflect.DeepEqual(searchPrQuery, r.URL.Query()) {
				searchPrJson, _ := testDataFs.ReadFile("testdata/search_prs.json")
				w.WriteHeader(http.StatusOK)
				w.Write(searchPrJson)
				return
			}

			if reflect.DeepEqual(searchReviewsQuery, r.URL.Query()) {
				searchReviewJson, _ := testDataFs.ReadFile("testdata/search_reviews.json")
				w.WriteHeader(http.StatusOK)
				w.Write(searchReviewJson)
				return
			}

      t.Log("Fall through in request handler")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
		}))

		defer server.Close()

		client := NewGithubClient("token", server.URL)
		got, err := client.FetchContributions("mickjuice", "2024-08-01")
		want := PullRequestReviewOverview{
			PullRequests: []PullRequestResult{
				{Title: "First PR", Date: "2024-08-01T00:00:00Z"},
				{Title: "Second PR", Date: "2024-08-03T00:00:00Z"},
			},
      Reviews: []PullRequestResult{
				{Title: "First Reviewed PR", Date: "2024-08-01T00:00:00Z"},
				{Title: "Another Review", Date: "2024-08-03T00:00:00Z"},

      },
		}

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}

	})
}
