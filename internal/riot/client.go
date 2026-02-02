package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/spf13/viper"
)

type AccountInfo struct {
	Name  string
	Tag   string
	Level int
	Rank  string
}

type SummonerResponse struct {
	PUUID         string `json:"puuid"`
	ProfileIconID int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
}

type LeagueEntry struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func FetchAccount(name, tag string) (*AccountInfo, error) {
	apiKey := viper.GetString("riot_api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("riot_api_key not found in config")
	}

	client := golio.NewClient(apiKey, golio.WithRegion(api.RegionBrasil))
	account, err := client.Riot.Account.GetByRiotID(name, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	level, err := fetchSummonerLevel(apiKey, account.Puuid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch summoner level: %w", err)
	}

	rank := fetchRank(apiKey, account.Puuid)

	return &AccountInfo{
		Name:  account.GameName,
		Tag:   account.TagLine,
		Level: level,
		Rank:  rank,
	}, nil
}

func fetchSummonerLevel(apiKey, puuid string) (int, error) {
	endpoint := fmt.Sprintf("https://br1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s", puuid)

	var summoner SummonerResponse
	if err := makeRiotAPIRequest(apiKey, endpoint, &summoner); err != nil {
		return 0, err
	}

	return summoner.SummonerLevel, nil
}

func fetchRank(apiKey, puuid string) string {
	encodedPUUID := url.QueryEscape(puuid)
	endpoint := fmt.Sprintf("https://br1.api.riotgames.com/lol/league/v4/entries/by-puuid/%s", encodedPUUID)

	var entries []LeagueEntry
	if err := makeRiotAPIRequest(apiKey, endpoint, &entries); err != nil {
		return "Unranked"
	}

	for _, entry := range entries {
		if entry.QueueType == "RANKED_SOLO_5x5" {
			return fmt.Sprintf("%s %s (%d LP)", entry.Tier, entry.Rank, entry.LeaguePoints)
		}
	}

	for _, entry := range entries {
		if entry.QueueType == "RANKED_FLEX_SR" {
			return fmt.Sprintf("%s %s (%d LP) [Flex]", entry.Tier, entry.Rank, entry.LeaguePoints)
		}
	}

	return "Unranked"
}

func makeRiotAPIRequest(apiKey, endpoint string, result interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}
