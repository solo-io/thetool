package addon

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
)

type EnableDisable struct {
}

func (e EnableDisable) configure(a *Addon) {
	var defaultSelection string
	if a.Enable {
		defaultSelection = "enable"
	} else {
		defaultSelection = "disable"
	}
	var question = []*survey.Question{
		{
			Name: "status",
			Prompt: &survey.Select{
				Message: "Enable addon status for " + a.Name,
				Options: []string{"enable", "disable"},
				Default: defaultSelection,
			},
			Validate: survey.Required,
		},
	}

	answer := struct {
		Status string
	}{}

	err := survey.Ask(question, &answer)
	if err != nil {
		fmt.Println("Unable to configure addon", a.Name, err)
	}

	a.Enable = answer.Status == "enable"
}

const (
	status     = "status"
	disable    = "disable"
	statsd     = "statsd"
	prometheus = "prometheus"
	all        = "all"
)

var (
	metricsStatus = map[string]string{
		disable:    "disable",
		statsd:     "use existing statsd",
		prometheus: "use existing prometheus",
		all:        "install everything",
	}
)

type MetricsConfigurator struct{}

func (m MetricsConfigurator) configure(a *Addon) {
	newStatus := askStatus(metricsStatus, a, []string{disable, statsd, prometheus, all})
	a.Enable = newStatus != disable
	a.Configuration[status] = newStatus
	if newStatus == statsd {
		askStatsdAddress(a)
	}
}

var (
	tracingStatus = map[string]string{
		disable:     "disable",
		"configure": "use existing jaeger",
		"install":   "install jaeger",
	}
)

type TracingConfigurator struct{}

func (t TracingConfigurator) configure(a *Addon) {
	newStatus := askStatus(tracingStatus, a, []string{disable, "configure", "install"})
	a.Enable = newStatus != disable
	a.Configuration[status] = newStatus
	if newStatus == "configure" {
		askJaegerAddress(a)
	}
}

func askStatus(m map[string]string, a *Addon, optionOrder []string) string {
	defaultSelection, ok := m[a.Configuration[status]]
	if !ok {
		defaultSelection = disable
	}
	options := make([]string, len(optionOrder))
	for i, o := range optionOrder {
		options[i] = m[o]
	}

	prompt := &survey.Select{
		Message: "Status for " + a.Name,
		Options: options,
		Default: defaultSelection,
	}

	var answer string
	err := survey.AskOne(prompt, &answer, survey.Required)
	if err != nil {
		fmt.Println("Unable to configure addon", a.Name, err)
		return a.Configuration[status]
	}

	for k, v := range m {
		if v == answer {
			return k
		}
	}
	return disable
}

func askStatsdAddress(a *Addon) {
	var questions = []*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Please enter the host for Statsd server",
			},
			Validate: survey.Required,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Please enter the port for Statsd server",
			},
			Validate: validatePort,
		},
	}

	answers := struct {
		Host string
		Port int
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println("Unable to get statsd address", err)
	}
	a.Configuration["statsd_host"] = answers.Host
	a.Configuration["statsd_port"] = strconv.Itoa(answers.Port)
}

func askJaegerAddress(a *Addon) {
	var questions = []*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Please enter the host for Jaeger server",
			},
			Validate: survey.Required,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Please enter the port for Jaeger server",
			},
			Validate: validatePort,
		},
	}

	answers := struct {
		Host string
		Port int
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println("Unable to get Jaeger address", err)
	}
	a.Configuration["jaeger_host"] = answers.Host
	a.Configuration["jaeger_port"] = strconv.Itoa(answers.Port)
}

func validatePort(val interface{}) error {
	if val == "" {
		return errors.New("port can't be empty")
	}

	port, err := strconv.Atoi(val.(string))
	if err != nil {
		return errors.Wrap(err, "port needs to be a number")
	}
	if port < 0 || port > 65535 {
		return errors.New("port not within a valid range")
	}
	return nil
}
