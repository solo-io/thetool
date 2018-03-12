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

type Addon struct {
	Name       string `json:"name"`
	Repository string `json:"repository,omitempty"`
	Commit     string `json:"commit,omitempty"`
	Image      string `json:"dockerImage,omitempty"`
	Tag        string `json:"dockerTag,omitempty"`
	Enable     bool   `json:"enable"`
	ConfigOnly *bool  `json:"configOnly,omitempty"`
}

var DefaultAddons = []*Addon{
	newGlooAddon("gloo-function-discovery",
		"https://github.com/solo-io/gloo-function-discovery.git",
		"644fefd36ce319638b8f4f5bab0ee20fb5a9f94c"),
	newGlooAddon("gloo-ingress-controller",
		"https://github.com/solo-io/gloo-ingress-controller.git",
		"90f2b216178ce58fe7dc9e1049e91d37f9a234fe"),
	newGlooAddon("gloo-k8s-service-discovery",
		"https://github.com/solo-io/gloo-k8s-service-discovery.git",
		"12b4753e52f6c7ab0d431a30b3f71f0b2caa5ff0"),
	newNonGlooAddon("statsd-exporter", "prom/statsd-exporter", "latest"),
	newNonGlooAddon("grafana", "grafana/grafana", "4.2.0"),
	newNonGlooAddon("prometheus", "quay.io/coreos/prometheus", "latest"),
	newNonGlooAddon("kube-state-metrics", "gcr.io/google_containers/kube-state-metrics", "v0.5.0"),
	newNonGlooAddon("jaeger", "jaegertracing/all-in-one", "latest"),
}

func newGlooAddon(name, repo, hash string) *Addon {
	return &Addon{
		Name:       name,
		Repository: repo,
		Commit:     hash,
		Enable:     true,
	}
}

func newNonGlooAddon(name, image, tag string) *Addon {
	return &Addon{
		Name:   name,
		Image:  image,
		Tag:    tag,
		Enable: false,
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
	if s.ConfigOnly != nil {
		fmt.Fprintf(b, "%-12s: %v\n", "Configure Only", *s.ConfigOnly)
	}

	return b.String()
}

func Init() error {
	return save(addonFilename, DefaultAddons)
}

func List() ([]*Addon, error) {
	return load(addonFilename)
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
