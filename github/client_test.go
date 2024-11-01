package github

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

//go:embed testdata/search_repo.json
var searchRepoJson string

func TestClient(t *testing.T) {
	repoPath := "/search/repositories"
	t.Run("FetchRepos returns list of repo names", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != repoPath {
				t.Errorf("Expected path %s, Got path %s", repoPath, r.URL.Path)
			}

			if authHeader := r.Header.Get("Authorization"); authHeader != "Bearer token" {
				t.Errorf("Expected Header: %q, got %q", "Bearer token", authHeader)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(searchRepoJson))
		}))
		defer server.Close()

		want := []string{"hey", "hi"}
		got, err := FetchRepos(server.URL, "token")

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
