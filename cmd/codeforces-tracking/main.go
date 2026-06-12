package codeforcestracking

import (
	"log"

	"cfvgtuai/internal/app"
	"cfvgtuai/internal/config"
)

func main() {
	cfg, err := config.Load("configs/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app.Run(cfg)
}
