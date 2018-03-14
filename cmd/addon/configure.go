package addon

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

func configureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "configure",
		Short:     "configure add-ons",
		ValidArgs: addonNames(),
		Args:      cobra.OnlyValidArgs,
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 1 {
				runConfigure(args[0])
			} else {
				name, err := askAddonName()
				if err != nil {
					fmt.Println("Unable to get addon to configure:", err)
					return
				}
				runConfigure(name)
			}
		},
	}
	return cmd
}

func runConfigure(name string) {
	addons, err := List()
	if err != nil {
		fmt.Println("Unable to get list of addons.")
		return
	}
	for _, a := range addons {
		if a.Name == name {
			configurator, ok := configuratorMap[a.Name]
			if !ok {
				fmt.Println("No configurator set for", a.Name)
				return
			}
			configurator.configure(a)
			if err := save(addonFilename, addons); err != nil {
				fmt.Println("Unable to update list of addons", err)
			}
			return
		}
	}
	fmt.Println("Unable to find addon named", name)
}

func askAddonName() (string, error) {
	var question = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Select{
				Message: "Select addon to configure",
				Options: addonNames(),
			},
			Validate: survey.Required,
		},
	}

	answer := struct {
		Name string
	}{}

	err := survey.Ask(question, &answer)
	if err != nil {
		return "", err
	}
	return answer.Name, nil
}
