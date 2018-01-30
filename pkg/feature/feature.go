package feature

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"
)

type Feature struct {
	Name       string `json:"name"`
	Version    string `json:"commit"`
	Repository string `json:"repository"`
	Enabled    bool   `json:"enabled"`
}

func ListDefaultFeatures() []Feature {
	// hard coded list of all known features
	return []Feature{
		Feature{
			Name:       "squash",
			Version:    "9397617b238cc4f17a0a3f0dc24194baf506ac97",
			Repository: "https://github.com/axhixh/envoy-squash.git",
			Enabled:    true,
		},
		Feature{
			Name:       "echo",
			Version:    "37a53fefe0a267fe3f4704c35e3721a4b6032f2a",
			Repository: "https://github.com/axhixh/echo.git",
			Enabled:    true,
		},
	}
}

func Save(features []Feature, w io.Writer) error {
	b, err := json.MarshalIndent(featureFile{
		Date:        time.Now(),
		GeneratedBy: "thetool",
		Features:    features,
	}, "", " ")
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

func SaveToFile(features []Feature, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return Save(features, f)
}

func Load(r io.Reader) ([]Feature, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 32*1024))
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	ff := &featureFile{}
	err = json.Unmarshal(buf.Bytes(), ff)
	if err != nil {
		return nil, err
	}
	return ff.Features, nil
}

func LoadFromFile(filename string) ([]Feature, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

type featureFile struct {
	Date        time.Time `json:"date"`
	GeneratedBy string    `json:"generatedBy"`
	Features    []Feature `json:"features"`
}
