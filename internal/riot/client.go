package riot

import (
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/spf13/viper"
)

type AccountInfo struct{
	Name 	string
	Tag		string
	Level	int
	Rank	string
	WinRate	string
}

func FetchAccount(name, tag string) (*AccountInfo, error) {
	apiKey := viper.GetString("riot_api_key")

	if apiKey == "" {
		return nil, fmt.Errorf("Api Key not found. Check your config.yaml")
	}
	client := golio.NewClient(apiKey, golio.WithRegion(api.RegionBrasil))

	account, err := client.Riot.Account.GetByRiotID(name, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	summoner, err := client.Riot.LoL.Summoner.GetByPUUID(account.Puuid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch summoner: %w", err)
	}

	entries, err := client.Riot.LoL.League.ListBySummoner(summoner.ID)
	rankString := "Unranked"

	for _, entry := range entries {
		if entry.QueueType == "RANKED_SOLO_5x5" {
			rankString = fmt.Sprintf("%s %s (%d LP)", entry.Tier, entry.Rank, entry.LeaguePoints)
		}
	}

	return &AccountInfo{
		Name:	account.GameName,
		Tag:	account.TagLine,
		Level:	summoner.SummonerLevel,
		Rank:	rankString,
	}, nil
}
