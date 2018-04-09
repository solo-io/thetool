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
	// EnvoyRepoUser - https://github.com/{EnvoyRepoUser}/envoy/
	EnvoyRepoUser = "envoyproxy"
	// EnvoyBuilderHash
	EnvoyBuilderHash = "6153d9787cb894c2dd6b17a1539eaeba88ae15d79f66f63eec0f4713436d74f0"

	// EnvoyCommonHash
	EnvoyCommonHash = "771b89c20a7a6f8edf3ebe3df2358f0e07e7edcd"

	// GlooHash is the commit hash of the version of Gloo used
	GlooHash = "2246f0e8e3e8739e0f2659ff114eb83e35ddd19d"
	// GlooRepo is the repository URL for Gloo
	GlooRepo = "https://github.com/solo-io/gloo.git"

	// DockerUser is the default Docker registry user used for publishing the images
	DockerUser = "soloio"
	// ConfigFile is the name of the configuraiton file
	ConfigFile = "thetool.json"
)

// Config contains the configuration used by thetool
type Config struct {
	EnvoyRepoUser    string `json:"envoyRepoUser"`
	EnvoyHash        string `json:"envoyHash"`
	EnvoyCommonHash  string `jsno:"envoyCommonHash"`
	EnvoyBuilderHash string `json:"envoyBuilderHash"`
	GlooHash         string `json:"glooHash"`
	GlooRepo         string `json:"glooRepo"`
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
