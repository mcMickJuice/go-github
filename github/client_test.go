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
        "q": {"org:shipt segway in:name"},
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
}
