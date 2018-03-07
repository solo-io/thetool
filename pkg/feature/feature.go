package feature

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	//FeaturesFileName represents the filename for features
	FeaturesFileName = "features.json"
)

type ManifestFeature struct {
	Name     string   `json:"name"`
	GlooDir  string   `json:"gloo,omitempty"`
	EnvoyDir string   `json:"envoy,omitempty"`
	Enabled  *bool    `json:"enabled,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func LoadManifest(folder string) ([]ManifestFeature, error) {
	b, err := ioutil.ReadFile(filepath.Join(folder, "features.json"))
	if err != nil {
		return nil, err
	}
	mf := []ManifestFeature{}
	err = json.Unmarshal(b, &mf)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func ToFeatures(repo, hash string, mf []ManifestFeature) []Feature {
	features := make([]Feature, len(mf))
	for i, f := range mf {
		enabled := true
		if f.Enabled != nil {
			enabled = *f.Enabled
		}
		features[i] = Feature{
			Name:       f.Name,
			GlooDir:    f.GlooDir,
			EnvoyDir:   f.EnvoyDir,
			Repository: repo,
			Revision:   hash,
			Enabled:    enabled,
			Tags:       f.Tags,
		}
	}
	return features
}

type Feature struct {
	Name       string   `json:"name"`
	GlooDir    string   `json:"gloo,omitempty"`
	EnvoyDir   string   `json:"envoy,omitempty"`
	Repository string   `json:"repository"`
	Revision   string   `json:"revision"`
	Enabled    bool     `json:"enabled"`
	Tags       []string `json:"tags,omitempty"`
}

type FeatureStore interface {
	Init() error
	AddAll([]Feature) error
	Update(Feature) error
	List() ([]Feature, error)
	RemoveForRepo(repoURL string) error
}
type FileFeatureStore struct {
	Filename string
}

func (f *FileFeatureStore) Init() error {
	return f.save([]Feature{})
}
func (f *FileFeatureStore) AddOrUpdateAll(features []Feature) error {
	if len(features) == 0 {
		// nothing to do
		return nil
	}
	existing, err := f.List()
	if err != nil {
		return err
	}

	// remove any existing features for this repo
	filtered := removeRepo(features[0].Repository, existing)
	var updated []Feature
	for _, feature := range features {
		for _, e := range filtered {
			if e.Name == feature.Name {
				return fmt.Errorf("feature %s already exists from another repository %s", feature.Name, e.Repository)
			}
		}
		updated = append(updated, feature)
	}
	return f.save(updated)
}

func (f *FileFeatureStore) Update(feature Feature) error {
	existing, err := f.List()
	if err != nil {
		return err
	}

	updated := make([]Feature, len(existing))
	for i, e := range existing {
		if e.Name == feature.Name {
			updated[i] = feature
		} else {
			updated[i] = e
		}
	}
	return f.save(updated)
}

func (f *FileFeatureStore) RemoveForRepo(repo string) error {
	existing, err := f.List()
	if err != nil {
		return err
	}

	updated := removeRepo(repo, existing)
	return f.save(updated)
}

func removeRepo(repo string, features []Feature) []Feature {
	updated := []Feature{}
	for _, e := range features {
		if e.Repository != repo {
			updated = append(updated, e)
		}
	}
	return updated
}

func (f *FileFeatureStore) save(features []Feature) error {
	b, err := json.MarshalIndent(featureFile{
		Date:        time.Now(),
		GeneratedBy: "thetool",
		Features:    features,
	}, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.Filename, b, 0644)
}

func (f *FileFeatureStore) List() ([]Feature, error) {
	b, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		return nil, err
	}
	ff := &featureFile{}
	err = json.Unmarshal(b, ff)
	if err != nil {
		return nil, err
	}
	return ff.Features, nil
}

type featureFile struct {
	Date        time.Time `json:"date"`
	GeneratedBy string    `json:"generatedBy"`
	Features    []Feature `json:"features"`
}
