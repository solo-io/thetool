package downloader

import "testing"

func TestRepoDir(t *testing.T) {
	cases := [][]string{
		{"git@github.com:solo-io/envoy-lambda.git", "envoy-lambda"},
		{"ssh://test@bitbucket.org:apple/ball", "ball"},
	}
	for _, c := range cases {
		out := repoDir(c[0])
		if c[1] != out {
			t.Errorf("expected %s got %q", c[1], out)
		}
	}
}
