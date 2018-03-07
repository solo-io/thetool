package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const (
	// WorkDir is the directory used by thetool
	WorkDir = "repositories"

	// EnvoyHash is the commit hash of the version of Envoy used
	EnvoyHash = "f79a62b7cc9ca55d20104379ee0576617630cdaa"
	// EnvoyBuilderHash
	EnvoyBuilderHash = "6153d9787cb894c2dd6b17a1539eaeba88ae15d79f66f63eec0f4713436d74f0"

	// GlooHash is the commit hash of the version of Gloo used
	GlooHash = "cf08737718cf62bf597f88aa2068c6f6b28b9992"
	// GlooRepo is the repository URL for Gloo
	GlooRepo = "https://github.com/solo-io/gloo.git"

	//GlooChartHash is the commit hash of the Gloo chart used
	GlooChartHash = "41cf74dabb6ee82752ed3887ba62c609047c9277"
	//GlooChartRepo is the repository URL for Gloo chart
	GlooChartRepo = "https://github.com/solo-io/gloo-install.git"

	// DockerUser is the default Docker registry user used for publishing the images
	DockerUser = "soloio"
	// ConfigFile is the name of the configuraiton file
	ConfigFile = "thetool.json"
)

// Config contains the configuration used by thetool
type Config struct {
	WorkDir          string `json:"workDir"`
	EnvoyHash        string `json:"envoyHash"`
	EnvoyBuilderHash string `json:"envoyBuilderHash"`
	GlooHash         string `json:"glooHash"`
	GlooRepo         string `json:"glooRepo"`
	GlooChartHash    string `json:"glooChartHash"`
	GlooChartRepo    string `json:"glooChartRepo"`
	DockerUser       string `json:"dockerUser,omitempty"`
}

// Save the current configuration used by thetool to a file
func (c *Config) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return saveToWriter(c, f)
}

func saveToWriter(c *Config, w io.Writer) error {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// Load the configuration for thetool from a file
func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return loadFromReader(f)
}

func loadFromReader(r io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = json.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
