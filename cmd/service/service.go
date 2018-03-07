package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

const (
	serviceFilename = "services.json"
)

type Service struct {
	Name       string `json:"name"`
	Repository string `json:"repository,omitifempty"`
	Commit     string `json:"commit,omitifempty"`
	Image      string `json:"dockerImage,omitifempty"`
	Tag        string `json:"dockerTag,omitifempty"`
	Enable     bool   `json:"enable"`
	Install    bool   `json:"install"`
}

var DefaultServices = []*Service{
	newGlooService("gloo-function-discovery",
		"https://github.com/solo-io/gloo-function-discovery.git",
		"644fefd36ce319638b8f4f5bab0ee20fb5a9f94c"),
	newGlooService("gloo-ingress",
		"https://github.com/solo-io/gloo-ingress.git",
		"99184ba6f4f35f8cfc416538b461deefcd6748bb"),
	newGlooService("gloo-k8s-service-discovery",
		"https://github.com/solo-io/gloo-k8s-service-discovery.git",
		"12b4753e52f6c7ab0d431a30b3f71f0b2caa5ff0"),
	newNonGlooService("grafana", "grafana/grafana", "4.2.0"),
	//newNonGlooService("grafana-dashboard", "giantswarm/tiny-tools", ""),
	newNonGlooService("prometheus", "quay.io/coreos/prometheus", "latest"),
	newNonGlooService("kube-state-metrics", "gcr.io/google_containers/kube-state-metrics", "v0.5.0"),
}

func newGlooService(name, repo, hash string) *Service {
	return &Service{
		Name:       name,
		Repository: repo,
		Commit:     hash,
		Enable:     true,
		Install:    true,
	}
}

func newNonGlooService(name, image, tag string) *Service {
	return &Service{
		Name:    name,
		Image:   image,
		Tag:     tag,
		Enable:  true,
		Install: true,
	}
}

func (s *Service) SafeName() string {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	return replacer.Replace(s.Name)
}
func (s *Service) ImageTag() string {
	if s.Tag != "" {
		return s.Tag
	}
	if len(s.Commit) >= 7 {
		return s.Commit[:7]
	}
	return ""
}
func (s *Service) String() string {
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
	fmt.Fprintf(b, "%-12s: %v\n", "Install", s.Install)

	return b.String()
}

func Init() error {
	return save(serviceFilename, DefaultServices)
}

func List() ([]*Service, error) {
	return load(serviceFilename)
}

// save and load; move to it to pkg/service?
type serviceFile struct {
	Date        time.Time  `json:"date"`
	GeneratedBy string     `json:"generatedBy"`
	Services    []*Service `json:"services"`
}

func save(filename string, services []*Service) error {
	b, err := json.MarshalIndent(serviceFile{
		Date:        time.Now(),
		GeneratedBy: "thetool",
		Services:    services,
	}, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, 0644)
}

func load(filename string) ([]*Service, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	sf := &serviceFile{}
	err = json.Unmarshal(b, sf)
	if err != nil {
		return nil, err
	}
	return sf.Services, nil
}
