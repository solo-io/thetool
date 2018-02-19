package glue

import "testing"

func TestUpdateDep(t *testing.T) {
	plugins := []GlooPlugin{
		GlooPlugin{
			Package:  "github.com/solo-io/gloo-plugins/aws",
			Revision: "aa23",
		},
	}
	updateDep(plugins, "test.txt")
}

func TestGetPackage(t *testing.T) {
	tests := [][]string{
		{"https://github.com/solo-io/gloo-plugins.git", "github.com/solo-io/gloo-plugins"},
		{"git@solo.io/test/as", "solo.io/test/as"},
		{"git@solo.io/test2/chori.git", "solo.io/test2/chori"},
	}

	for _, entry := range tests {
		if getPackage(entry[0]) != entry[1] {
			t.Errorf("expected %s got %s", entry[1], getPackage(entry[0]))
		}
	}
}
