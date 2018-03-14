package addon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

const (
	addonFilename = "addons.json"
)

type configurator interface {
	configure(*Addon)
}
type Addon struct {
	Name          string            `json:"name"`
	Repository    string            `json:"repository,omitempty"`
	Commit        string            `json:"commit,omitempty"`
	Image         string            `json:"dockerImage,omitempty"`
	Tag           string            `json:"dockerTag,omitempty"`
	Enable        bool              `json:"enable"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

var configuratorMap map[string]configurator = make(map[string]configurator)
var DefaultAddons = []*Addon{
	newGlooAddon("gloo-function-discovery",
		"https://github.com/solo-io/gloo-function-discovery.git",
		"51580349b03ece51bcd50997969584d143f6422b"),
	newGlooAddon("gloo-ingress-controller",
		"https://github.com/solo-io/gloo-ingress-controller.git",
		"90f2b216178ce58fe7dc9e1049e91d37f9a234fe"),
	newGlooAddon("gloo-k8s-service-discovery",
		"https://github.com/solo-io/gloo-k8s-service-discovery.git",
		"12b4753e52f6c7ab0d431a30b3f71f0b2caa5ff0"),
	tracingAddon(),
	metricsAddon(),
}

func newGlooAddon(name, repo, hash string) *Addon {
	configuratorMap[name] = EnableDisable{}
	return &Addon{
		Name:       name,
		Repository: repo,
		Commit:     hash,
		Enable:     true,
	}
}

func tracingAddon() *Addon {
	name := "opentracing"
	configuratorMap[name] = TracingConfigurator{}
	return &Addon{
		Name:   name,
		Enable: false,
		Configuration: map[string]string{
			"jaeger": "jaegertracing/all-in-one:latest",
			"status": "install",
		},
	}
}

func metricsAddon() *Addon {
	name := "metrics"
	configuratorMap[name] = MetricsConfigurator{}
	return &Addon{
		Name:   name,
		Enable: false,
		Configuration: map[string]string{
			"statsd-exporter": "prom/statsd-exporter:latest",
			"grafana":         "grafana/grafana:4.2.0",
			"prometheus":      "quay.io/coreos/prometheus:latest",
			"status":          "install",
		},
	}
}

func (s *Addon) SafeName() string {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	return replacer.Replace(s.Name)
}
func (s *Addon) ImageTag() string {
	if s.Tag != "" {
		return s.Tag
	}
	if len(s.Commit) >= 7 {
		return s.Commit[:7]
	}
	return ""
}
func (s *Addon) String() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%-12s: %s\n", "Name", s.Name)
	if s.Repository != "" {
		fmt.Fprintf(b, "%-12s: %s\n", "Repository", s.Repository)
		fmt.Fprintf(b, "%-12s: %s\n", "Commit", s.Commit)
	}
	if s.Image != "" {
		fmt.Fprintf(b, "%-12s: %s\n", "Image", s.Image)
		fmt.Fprintf(b, "%-12s: %s\n", "Tag", s.Tag)
	}
	fmt.Fprintf(b, "%-12s: %v\n", "Enable", s.Enable)
	if s.Enable && s.Configuration != nil {
		status, ok := s.Configuration["status"]
		if ok {
			fmt.Fprintf(b, "%-12s: %s\n", "Status", status)
		}
	}
	return b.String()
}

func Init() error {
	return save(addonFilename, DefaultAddons)
}

func List() ([]*Addon, error) {
	return load(addonFilename)
}

func addonNames() []string {
	addons, err := List()
	if err != nil {
		return []string{}
	}
	names := make([]string, len(addons))
	for i, a := range addons {
		names[i] = a.Name
	}
	return names
}

// save and load; move to it to pkg/addon?
type addonFile struct {
	Date        time.Time `json:"date"`
	GeneratedBy string    `json:"generatedBy"`
	Addons      []*Addon  `json:"addons"`
}

func save(filename string, addons []*Addon) error {
	b, err := json.MarshalIndent(addonFile{
		Date:        time.Now(),
		GeneratedBy: "thetool",
		Addons:      addons,
	}, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, 0644)
}

func load(filename string) ([]*Addon, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	sf := &addonFile{}
	err = json.Unmarshal(b, sf)
	if err != nil {
		return nil, err
	}
	return sf.Addons, nil
}
