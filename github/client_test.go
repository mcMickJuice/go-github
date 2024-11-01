package github

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient(t *testing.T) {
	repoPath := "/orgs/shipt/repos"
	t.Run("FetchRepos returns list of repo names", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != repoPath {
				t.Errorf("Expected path %s, Got path %s", repoPath, r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": 1, "name": "hey"},{"id": 2,"name": "hi"}]`))
		}))
		defer server.Close()

		want := []string{"hey", "hi"}
		got, err := FetchRepos(server.URL)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
