package feature

import "testing"

func TestManifest(t *testing.T) {
	mf, err := LoadManifest("testdata")
	if err != nil {
		t.Error("failed loading manifest", err)
	}
	features := ToFeatures("repo", "hash", mf)
	if len(features) != 2 {
		t.Error("expected 2 features but got ", len(features))
	}

	if !features[0].Enabled {
		t.Error("expected first feature to be enabled")
	}

	if features[1].Enabled {
		t.Error("expected second feature to be disabled")
	}

	if len(features[1].Tags) != 1 {
		t.Error("expected there to be a tag")
	}
}
