package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/somnifobia/league_account_clitool/internal/riot"
	"github.com/somnifobia/league_account_clitool/internal/store"
)

var addCmd = &cobra.Command{
	Use:   "add [gameName] [tagLine]",
	Short: "Add an account",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		tag := args[1]

		fmt.Printf("Searching data for %s#%s...\n", name, tag)

		riotInfo, err := riot.FetchAccount(name, tag)
		if err != nil {
			fmt.Printf("Error fetching account: %v\n", err)
			return
		}

		fmt.Printf("Account found: Level %d | Elo: %s\n", riotInfo.Level, riotInfo.Rank)
		fmt.Print("Type the actual BlueEssence amount: ")

		var beInput string
		fmt.Scanln(&beInput)
		be, _ := strconv.Atoi(beInput)

		acc := store.Account{
			Name:		riotInfo.Name,
			Tag:		riotInfo.Tag,
			Level:		riotInfo.Level,
			Rank:		riotInfo.Rank,
			BlueEssence: be,
		}

		if err := store.AddAccount(acc); err != nil {
			fmt.Println("Error adding account: ", err)
			return
		}

		fmt.Println("Account added successfully!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
