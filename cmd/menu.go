package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var menuOptions = []string{
	"Add account",
	"List accounts",
	"Delete account",
	"Exit",
}

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Open interactive menu",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			var selection string
			prompt := &survey.Select{
				Message: "What you want to do?",
				Options: menuOptions,
			}
			survey.AskOne(prompt, &selection)

			switch selection {
			case "Add account":
				interactiveAdd()
			case "List accounts":
				// Chama o comando 'list' que já existe
				listCmd.Run(cmd, []string{})
				waitEnter()
			case "Delete account":
				interactiveDelete()
			case "Exit":
				fmt.Println("bye")
				os.Exit(0)
			}
		}
	},
}

func interactiveAdd() {
	var answers struct {
		Name string
		Tag  string
	}

	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{Message: "Type summoner name:"},
			Validate: survey.Required,
		},
		{
			Name: "tag",
			Prompt: &survey.Input{Message: "Type tag(ex: BR1):"},
			Validate: survey.Required,
		},
	}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println("Canceled.")
		return
	}

	fmt.Println("\n--- Starting Add Process ---")
	addCmd.Run(nil, []string{answers.Name, answers.Tag})
	waitEnter()
}

func interactiveDelete() {
	var nameToDelete string
	prompt := &survey.Input{
		Message: "Type account name to delete:",
	}
	survey.AskOne(prompt, &nameToDelete)

	deleteCmd.Run(nil, []string{nameToDelete})
	waitEnter()
}

func waitEnter() {
	fmt.Print("\nPress Enter to continue...")
	fmt.Scanln()
}

func init() {
	rootCmd.AddCommand(menuCmd)
}
