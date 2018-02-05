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
	EnvoyHash = "29989a38c017d3be5aa3c735a797fcf58b754fe5"
	// GlueHash is the commit hash of the version of Glue used
	GlueHash = "5309cb36385555b7c2d5278fc230b2b27d8a0787"
	// DockerUser is the default Docker registry user used for publishing the images
	DockerUser = "solo-io"
	// ConfigFile is the name of the configuraiton file
	ConfigFile = "thetool.json"
)

// Config contains the configuration used by thetool
type Config struct {
	WorkDir    string `json:"workDir"`
	EnvoyHash  string `json:"envoyHash"`
	GlueHash   string `json:"glueHash"`
	DockerUser string `json:"dockerUser,omitempty"`
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
