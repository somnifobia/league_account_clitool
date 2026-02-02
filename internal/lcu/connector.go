package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type Wallet struct {
	LolBlueEssence int `json:"lol_blue_essence"`
	RP             int `json:"RP"`
}

func GetWallet() (int, error) {
	port, password, err := getLockfileData()
	if err != nil {
		return 0, fmt.Errorf("League isnt open: %v", err)
	}

	url := fmt.Sprintf("https://127.0.0.1:%s/lol-inventory/v1/wallet?inventoryType=CHAMPION", port)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("riot", password))

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var wallet Wallet
	if err := json.Unmarshal(body, &wallet); err != nil {
		return 0, err
	}

	return wallet.LolBlueEssence, nil
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
			cmd, _ := p.Cmdline()
			exePath, _ := p.Exe()
			dir := filepath.Dir(exePath)
			lockfile := filepath.Join(dir, "lockfile")

			data, err := os.ReadFile(lockfile)
			if err != nil {
				return parseCmdLine(cmd)
			}
			parts := strings.Split(string(data), ":")
			if len(parts) >= 3 {
				return parts[2], parts[3], nil
			}
		}
	}
	return "", "", fmt.Errorf("Process not found")
}

func parseCmdLine(cmd string) (string, string, error) {
	return "", "", fmt.Errorf("Command line parsing not implemented")
}
