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
	// EnvoyBuilderHash
	EnvoyBuilderHash = "6153d9787cb894c2dd6b17a1539eaeba88ae15d79f66f63eec0f4713436d74f0"

	// GlooHash is the commit hash of the version of Gloo used
	GlooHash = "c88c90c332e5528a070a1c800bc65b2c39f8ca24"
	// GlooRepo is the repository URL for Gloo
	GlooRepo = "https://github.com/solo-io/gloo.git"

	//GlooChartHash is the commit hash of the Gloo chart used
	GlooChartHash = "a2f12f82fd41d7b7eab91a2f825fee8f9fdf6ec5"
	//GlooChartRepo is the repository URL for Gloo chart
	GlooChartRepo = "https://github.com/solo-io/gloo-chart.git"

	GlooFuncDiscoveryHash = "644fefd36ce319638b8f4f5bab0ee20fb5a9f94c"
	GlooFuncDiscoveryRepo = "https://github.com/solo-io/gloo-function-discovery.git"

	GlooIngressHash = "99184ba6f4f35f8cfc416538b461deefcd6748bb"
	GlooIngressRepo = "https://github.com/solo-io/gloo-ingress.git"

	GlooK8SDiscoveryHash = "12b4753e52f6c7ab0d431a30b3f71f0b2caa5ff0"
	GlooK8SDiscvoeryRepo = "https://github.com/solo-io/gloo-k8s-service-discovery.git"

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
	GlooFuncDHash    string `json:"glooFuncDiscoveryHash"`
	GlooFuncDRepo    string `json:"glooFuncDiscoveryRepo"`
	GlooIngressHash  string `json:"glooIngressHash"`
	GlooIngressRepo  string `json:"glooIngressRepo"`
	GlooK8SDHash     string `json:"glooK8SDiscoveryHash"`
	GlooK8SDRepo     string `json:"glooK8SDiscoveryRepo"`
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
