/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/somnifobia/league_account_clitool/internal/store"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove account by name",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := store.RemoveAccount(name)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Account '%s' deleted successfully\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
