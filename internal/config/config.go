package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	TelegramToken      string        `json:"telegram_token"`
	TelegramChatID     int64         `json:"telegram_chat_id"`
	CodeforcesHandles  []string      `json:"codeforces_handles"`
	PollingInterval    time.Duration `json:"polling_interval_seconds"`
	CodeforcesAPIToken string        `json:"codeforces_api_token"`
	CodeforcesAPIKey   string        `json:"codeforces_api_key"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	cfg.PollingInterval *= time.Second

	return cfg, nil
}
