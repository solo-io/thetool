package addon

import (
	"fmt"

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

var (
	metricsStatus = map[string]string{
		"disable":    "disable",
		"statsd":     "use existing statsd",
		"prometheus": "use existing prometheus",
		"all":        "install everything",
	}
)

type MetricsConfigurator struct{}

func (m MetricsConfigurator) configure(a *Addon) {
	defaultSelection, ok := metricsStatus[a.Configuration["status"]]
	if !ok {
		defaultSelection = "disable"
	}
	var question = []*survey.Question{
		{
			Name: "status",
			Prompt: &survey.Select{
				Message: "Status for " + a.Name,
				// for now not generating this from metricsStatus map to preserver order
				Options: []string{"disable", "use existing statsd", "use existing prometheus", "install everything"},
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

	a.Enable = answer.Status != "disable"
	a.Configuration["status"] = selectionToStatus(answer.Status, metricsStatus)
}

func selectionToStatus(selection string, m map[string]string) string {
	for k, v := range m {
		if v == selection {
			return k
		}
	}
	return "disable"
}

var (
	tracingStatus = map[string]string{
		"disable":   "disable",
		"configure": "use existing jaeger",
		"install":   "install jaeger",
	}
)

type TracingConfigurator struct{}

func (t TracingConfigurator) configure(a *Addon) {
	defaultSelection, ok := tracingStatus[a.Configuration["status"]]
	if !ok {
		defaultSelection = "disable"
	}
	var question = []*survey.Question{
		{
			Name: "status",
			Prompt: &survey.Select{
				Message: "Status for " + a.Name,
				Options: []string{"disable", "use existing jaeger", "install jaeger"},
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

	a.Enable = answer.Status != "disable"
	a.Configuration["status"] = selectionToStatus(answer.Status, tracingStatus)
}
