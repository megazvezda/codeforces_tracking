package app

import (
	"log"
	"sort"
	"time"

	"cfvgtuai/internal/codeforces"
	"cfvgtuai/internal/config"
	"cfvgtuai/internal/telegram"
)

type Tracker struct {
	processedSubmissions map[int]bool
	wrongAttempts        map[string]int
}

func NewTracker() *Tracker {
	return &Tracker{
		processedSubmissions: make(map[int]bool),
		wrongAttempts:        make(map[string]int),
	}
}

func Run(cfg *config.Config) {
	cfClient := codeforces.NewClient()
	tgNotifier, err := telegram.NewNotifier(cfg.TelegramToken, cfg.TelegramChatID)
	if err != nil {
		log.Fatalf("failed to create telegram notifier: %v", err)
	}

	tracker := NewTracker()

	for {
		for _, handle := range cfg.CodeforcesHandles {
			submissions, err := cfClient.GetUserSubmissions(handle)
			if err != nil {
				log.Printf("failed to get submissions for %s: %v", handle, err)
				continue
			}

			sort.Slice(submissions, func(i, j int) bool {
				return submissions[i].CreationTimeSeconds < submissions[j].CreationTimeSeconds
			})

			for _, submission := range submissions {
				if tracker.processedSubmissions[submission.ID] {
					continue
				}

				problemID := submission.Problem.Name
				if submission.Verdict == "OK" {
					if !tracker.processedSubmissions[submission.ID] {
						wrongAttempts := tracker.wrongAttempts[problemID]
						tgNotifier.Notify(submission, handle, wrongAttempts)
						tracker.processedSubmissions[submission.ID] = true
						delete(tracker.wrongAttempts, problemID)
					}
				} else {
					tracker.wrongAttempts[problemID]++
				}
				tracker.processedSubmissions[submission.ID] = true
			}
		}
		time.Sleep(cfg.PollingInterval)
	}
}
