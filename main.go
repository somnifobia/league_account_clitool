package main

import(
	"os"
	"github.com/somnifobia/league_account_clitool/cmd"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "menu")
	}

	cmd.Execute()
}
