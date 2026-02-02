package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type WalletStore struct {
	IP int `json:"ip"`
	RP int `json:"rp"`
}

type WalletInventory struct {
	LolBlueEssence int `json:"lol_blue_essence"`
	RP             int `json:"RP"`
}

type LootMap map[string]LootItem

type LootItem struct {
	Asset string `json:"asset"`
	Count int    `json:"count"`
}

func GetWallet() (int, error) {
	port, password, err := getLockfileData()
	if err != nil {
		return 0, fmt.Errorf("client not found/open: %v", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}

	be, err := tryStoreWallet(client, port, password)
	if err == nil {
		return be, nil
	}

	be, err = tryInventoryWallet(client, port, password)
	if err == nil {
		return be, nil
	}

	be, err = tryLootMap(client, port, password)
	if err == nil {
		return be, nil
	}

	return 0, fmt.Errorf("failed to detect BE on all endpoints")
}

func tryStoreWallet(client *http.Client, port, password string) (int, error) {
	url := fmt.Sprintf("https://127.0.0.1:%s/lol-store/v1/wallet", port)
	body, err := doRequest(client, url, password)
	if err != nil {
		return 0, err
	}

	var w WalletStore
	if err := json.Unmarshal(body, &w); err != nil {
		return 0, err
	}
	return w.IP, nil
}

func tryInventoryWallet(client *http.Client, port, password string) (int, error) {
	params := url.Values{}
	params.Set("inventoryTypes", "[\"CHAMPION\"]")
	query := params.Encode()

	url := fmt.Sprintf("https://127.0.0.1:%s/lol-inventory/v1/wallet?%s", port, query)

	body, err := doRequest(client, url, password)
	if err != nil {
		return 0, err
	}

	var w WalletInventory
	if err := json.Unmarshal(body, &w); err != nil {
		return 0, err
	}
	return w.LolBlueEssence, nil
}

func tryLootMap(client *http.Client, port, password string) (int, error) {
	url := fmt.Sprintf("https://127.0.0.1:%s/lol-loot/v1/player-loot-map", port)
	body, err := doRequest(client, url, password)
	if err != nil {
		return 0, err
	}

	var m LootMap
	if err := json.Unmarshal(body, &m); err != nil {
		return 0, err
	}

	if item, ok := m["CURRENCY_champion"]; ok {
		return item.Count, nil
	}

	return 0, fmt.Errorf("currency not found in loot map")
}

func doRequest(client *http.Client, url, password string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("riot", password))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getLockfileData() (string, string, error) {
	processes, _ := process.Processes()
	for _, p := range processes {
		name, _ := p.Name()
		if name == "LeagueClientUx.exe" {
			exePath, _ := p.Exe()
			dir := filepath.Dir(exePath)
			lockfile := filepath.Join(dir, "lockfile")
			data, err := os.ReadFile(lockfile)
			if err != nil {
				return "", "", err
			}
			parts := strings.Split(string(data), ":")
			if len(parts) >= 3 {
				return parts[2], parts[3], nil
			}
		}
	}
	return "", "", fmt.Errorf("LeagueClientUx.exe not found")
}
