package gloo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/solo-io/thetool/pkg/feature"
)

func TestUpdateDep(t *testing.T) {
	plugins := []GlooPlugin{
		GlooPlugin{
			Package:  "github.com/solo-io/gloo-plugins/aws",
			Revision: "aa23",
		},
	}
	updateDep(plugins, "test.txt", "")
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

func TestIntallPlugins(t *testing.T) {
	tmp, err := ioutil.TempFile("", "thetool-test")
	if err != nil {
		t.Error("unable to create temporary file", err)
	}
	filename := tmp.Name()
	defer os.Remove(filename)

	plugins := []GlooPlugin{
		GlooPlugin{Package: "github.com/solo-io/gloo/pkg/plugins/aws", Revision: "ac23", Repository: "https://github.com/solo-io/gloo.git"},
		GlooPlugin{Package: "bitbucket.org/axhixh/gloo-plugins/magic", Revision: "23asc", Repository: "https://bitbucket.org/axhixh/gloo-plugins.git"},
	}

	if err := installPlugins(plugins, filename, installTemplate); err != nil {
		t.Error("unable to generate plugin install file", err)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("unable to read generated plugin install file", err)
	}

	for _, p := range plugins {
		if !bytes.Contains(data, []byte(fmt.Sprintf("_ \"%s\"", p.Package))) {
			t.Error("package not imported:", p.Package)
		}
	}
}

func TestToGlooPlugins(t *testing.T) {
	var cases = []struct {
		Features []feature.Feature
		Packages []string
	}{
		{ // empty
			[]feature.Feature{},
			[]string{},
		},
		{ // non gloo plugin
			[]feature.Feature{
				feature.Feature{Name: "test", EnvoyDir: "test/envoy", Revision: "1", Repository: "https://g.com/user/test.git"},
			},
			[]string{},
		},
		{
			[]feature.Feature{
				feature.Feature{Name: "aws", GlooDir: "pkg/plugins/aws", Revision: "232ac3d", Repository: "https://github.com/solo-io/gloo.git"},
			},
			[]string{"github.com/solo-io/gloo/pkg/plugins/aws"},
		},
	}
	for _, tc := range cases {
		plugins := toGlooPlugins(tc.Features)
		if len(plugins) != len(tc.Packages) {
			t.Errorf("expected %d plugins got %d", len(tc.Packages), len(plugins))
		}

		for i, p := range tc.Packages {
			if p != plugins[i].Package {
				t.Errorf("packages don't match; expected %s got %s", p, plugins[i].Package)
			}
		}
	}
}
