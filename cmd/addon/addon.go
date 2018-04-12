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
	Metrics       = "metrics"
	OpenTracing   = "opentracing"

	// configuration keys
	keyStatus = "status"
	keyEnable = "enable"
	keyGloo   = "gloo"
)

type configurator interface {
	configure(*Addon)
}
type Addon struct {
	Name          string                 `json:"name"`
	Configuration map[string]interface{} `json:"configuration"`
}

var configuratorMap map[string]configurator = make(map[string]configurator)
var DefaultAddons = []*Addon{
	newGlooAddon("function-discovery"),
	newGlooAddon("kube-ingress-controller"),
	newGlooAddon("kube-upstream-discovery"),
	tracingAddon(),
	metricsAddon(),
}

func newGlooAddon(name string) *Addon {
	configuratorMap[name] = EnableDisable{}
	return &Addon{
		Name:          name,
		Configuration: map[string]interface{}{keyEnable: true, keyGloo: true},
	}
}

func tracingAddon() *Addon {
	name := OpenTracing
	configuratorMap[name] = TracingConfigurator{}
	return &Addon{
		Name: name,
		Configuration: map[string]interface{}{
			"jaeger":  "jaegertracing/all-in-one:latest",
			keyStatus: "disable",
		},
	}
}

func metricsAddon() *Addon {
	name := Metrics
	configuratorMap[name] = MetricsConfigurator{}
	return &Addon{
		Name: name,
		Configuration: map[string]interface{}{
			"statsd_exporter": "prom/statsd-exporter:latest",
			keyStatus:         "disable",
		},
	}
}

func (s *Addon) SafeName() string {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	return replacer.Replace(s.Name)
}

func (s *Addon) IsGlooAddon() bool {
	return s.Configuration != nil && s.Configuration[keyGloo] == true
}

func (s *Addon) String() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%-12s: %s\n", "Name", s.Name)
	if s.Configuration != nil {
		status, ok := s.Configuration[keyStatus]
		if ok {
			fmt.Fprintf(b, "%-12s: %v\n", "Status", status)
		}

		enabled, ok := s.Configuration[keyEnable]
		if ok {
			fmt.Fprintf(b, "%-12s: %v\n", "Enable", enabled)
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

// Need a more generic solution
func InstallPrometheus() bool {
	addons, err := load(addonFilename)
	if err != nil {
		return false
	}
	for _, a := range addons {
		if a.Name == Metrics {
			status := a.Configuration[keyStatus]
			return "all" == status
		}
	}
	return false
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
