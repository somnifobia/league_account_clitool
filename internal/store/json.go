package store

import (
	"encoding/json"
	"os"
	"fmt"
)

type Account struct {
	Name		string	`json:"name"`
	Tag			string	`json:"tag"`
	Level		int		`json:"level"`
	Rank		string	`json:"rank"`
	BlueEssence	int		`json:"blue_essence"`
}

const dbFile = "accounts.json"

func AddAccount(acc Account) error {
	accounts, _ := ListAccounts()
	accounts = append(accounts, acc)
	return saveFile(accounts)
}

func ListAccounts() ([]Account, error) {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		return []Account{}, err
	}
	var accounts []Account
	err = json.Unmarshal(data, &accounts)
	return accounts, err
}

func RemoveAccount(nameToRemove string) error {
	accounts, err := ListAccounts()
	if err != nil {
		return err
	}

	newAccounts := []Account{}
	found := false

	for _, acc := range accounts {
		if acc.Name != nameToRemove {
			newAccounts = append(newAccounts, acc)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("Account '%s' not found", nameToRemove)
	}

	return saveFile(newAccounts)
}


func saveFile(accounts []Account) error {
	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbFile, data, 0644)
}
