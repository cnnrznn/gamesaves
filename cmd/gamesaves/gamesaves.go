package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/cnnrznn/gamesaves/pkg/store/googledrive"
)

type Config struct {
	Backups []Backup `json:"backups"`
}

type Backup struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
}

func main() {
	configFile := os.Getenv("GAMESAVES_CONFIG_FILE")
	if configFile == "" {
		configFile = "/etc/gamesaves/config.json"
	}

	tokenFile := os.Getenv("GAMESAVES_TOKEN_FILE")
	if tokenFile == "" {
		tokenFile = "/etc/gamesaves/token.json"
	}

	_, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
		return
	}

	var noToken bool
	bs, err := os.ReadFile(tokenFile)
	if err != nil {
		noToken = true
	}

	if noToken {
		token, err := authorize()
		if err != nil {
			log.Fatal("failed to authorize storage: %w", err)
			return
		}

		f, err := os.Create(tokenFile)
		if err != nil {
			log.Fatal("failed to create token file: %w", err)
			return
		}

		err = json.NewEncoder(f).Encode(token)
		if err != nil {
			log.Fatal("failed to encode token: %w", err)
			return
		}

		return
	}

	var token *oauth2.Token
	err = json.Unmarshal(bs, &token)
	if err != nil {
		log.Fatal(err)
		return
	}

	// TODO read args to determine operating mode. Download/Upload and game name
}

func loadConfig(fn string) (*Config, error) {
	bs, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var config = &Config{}

	err = json.Unmarshal(bs, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func authorize() (*oauth2.Token, error) {
	token, err := googledrive.Authorize()
	if err != nil {
		return nil, err
	}

	return token, nil
}
