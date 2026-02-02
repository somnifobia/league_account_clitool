# League Account Manager
CLI tool written in Go for managing League of Legends accounts locally.
It fetches public data (Rank, Level) via the Riot API and detects private data (Blue Essence) via the local LCU client.

## Features
- Add accounts by fetching Summoner Level and Rank via Riot API
- Automatically detect Blue Essence by connecting to the local League Client (LCU)
- Persist data locally using a JSON file
- Interactive menu (TUI) for easy navigation without memorizing commands

## Prerequisites
- Riot API Key  
A valid development key from https://developer.riotgames.com/

- League of Legends Client  
Must be running to enable Blue Essence auto-detection

## Installation
### From Source
Clone the repository:
git clone https://github.com/somnifobia/league_account_clitool
cd league_account_clitool

Install dependencies:
go mod tidy

Build the executable:
go build -o lol-manager.exe

## Configuration
Create a file named config.yaml in the root directory.

Add your Riot API key:
riot_api_key: "RGAPI-YOUR-KEY-HERE"

## Usage
Run the executable to open the interactive menu:
./lol-manager.exe

## CLI Commands
Add Account:
./lol-manager.exe add "SummonerName" "TagLine"

Example:
./lol-manager.exe add "Faker" "KR1"

List Accounts:
./lol-manager.exe list

Delete Account:
./lol-manager.exe delete "SummonerName"

## Troubleshooting
Api Key not found:
Ensure that config.yaml exists in the same directory as the executable.

Blue Essence Detection Failed:
Make sure the League Client is open and logged in.
If detection fails, the tool will prompt for manual input.

401 / 403 API Errors:
Your development API key may have expired (24-hour validity).
Generate a new key on the Riot Developer Portal.

## License
This project is licensed under the MIT License.
