package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/somnifobia/league_account_clitool/internal/store"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts saved",
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := store.ListAccounts()
		if err != nil {
			fmt.Println("Error reading from database", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Nome", "Tag", "Level", "Rank", "Blue Essence"})
		table.SetBorder(false)

		for _, acc := range accounts {
			table.Append([]string{
				acc.Name,
				acc.Tag,
				strconv.Itoa(acc.Level),
				acc.Rank,
				strconv.Itoa(acc.BlueEssence),
			})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
