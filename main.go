package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramBotToken    string   `json:"telegram_bot_token"`
	TelegramChatID      int64    `json:"telegram_chat_id"`
	CodeforcesHandles   []Handle `json:"codeforces_handles"`
	PollIntervalSeconds int      `json:"poll_interval_seconds"`
}

type Handle struct {
	Handle   string `json:"handle"`
	RealName string `json:"real_name"`
}

type CFResponse struct {
	Status string       `json:"status"`
	Result []Submission `json:"result"`
}

type Submission struct {
	ID                  int64   `json:"id"`
	ContestID           int     `json:"contestId"`
	CreationTimeSeconds int64   `json:"creationTimeSeconds"`
	Problem             Problem `json:"problem"`
	Author              Author  `json:"author"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict"`
	TimeConsumedMillis  int     `json:"timeConsumedMillis"`
	MemoryConsumedBytes int64   `json:"memoryConsumedBytes"`
	PassedTestCount     int     `json:"passedTestCount"`
}

type Problem struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Tags      []string `json:"tags"`
}

type Author struct {
	Members []Member `json:"members"`
}

type Member struct {
	Handle string `json:"handle"`
}

type CodeforcesClient struct {
	lastRequest time.Time
}

func NewCodeforcesClient() *CodeforcesClient {
	return &CodeforcesClient{
		lastRequest: time.Now().Add(-3 * time.Second),
	}
}

func (c *CodeforcesClient) GetUserSubmissions(handle string, count int) ([]Submission, error) {
	// Rate limiting: wait at least 2 seconds between requests
	elapsed := time.Since(c.lastRequest)
	if elapsed < 2*time.Second {
		time.Sleep(2*time.Second - elapsed)
	}

	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=%d", handle, count)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}
	defer resp.Body.Close()

	c.lastRequest = time.Now()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var cfResp CFResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if cfResp.Status != "OK" {
		return nil, fmt.Errorf("API returned status: %s", cfResp.Status)
	}

	return cfResp.Result, nil
}

type Tracker struct {
	processedSubmissions map[string]map[int64]bool // handle -> submissionID -> processed
	problemAttempts      map[string]map[string]int // handle -> problemKey -> attempts
}

func NewTracker() *Tracker {
	return &Tracker{
		processedSubmissions: make(map[string]map[int64]bool),
		problemAttempts:      make(map[string]map[string]int),
	}
}

func (t *Tracker) IsProcessed(handle string, submissionID int64) bool {
	if _, ok := t.processedSubmissions[handle]; !ok {
		t.processedSubmissions[handle] = make(map[int64]bool)
	}
	return t.processedSubmissions[handle][submissionID]
}

func (t *Tracker) MarkProcessed(handle string, submissionID int64) {
	if _, ok := t.processedSubmissions[handle]; !ok {
		t.processedSubmissions[handle] = make(map[int64]bool)
	}
	t.processedSubmissions[handle][submissionID] = true
}

func (t *Tracker) CountWrongAttempts(handle string, problemKey string, submissions []Submission) int {
	wrongCount := 0
	for _, sub := range submissions {
		subKey := fmt.Sprintf("%d%s", sub.Problem.ContestID, sub.Problem.Index)
		if subKey == problemKey && sub.Verdict != "OK" {
			wrongCount++
		}
	}
	return wrongCount
}

func formatNotification(handle, realName string, sub Submission, wrongAttempts int) string {
	// Detect if it's a gym contest (typically contestId >= 100000)
	// Gym contests use different URL format
	var problemURL, submissionURL string
	isGym := sub.ContestID >= 100000

	if isGym {
		problemURL = fmt.Sprintf("https://codeforces.com/gym/%d/problem/%s", sub.ContestID, sub.Problem.Index)
		submissionURL = fmt.Sprintf("https://codeforces.com/gym/%d/submission/%d", sub.ContestID, sub.ID)
	} else {
		problemURL = fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s", sub.ContestID, sub.Problem.Index)
		submissionURL = fmt.Sprintf("https://codeforces.com/contest/%d/submission/%d", sub.ContestID, sub.ID)
	}

	var attemptsText string
	if wrongAttempts == 0 {
		attemptsText = "without wrong submissions"
	} else if wrongAttempts == 1 {
		attemptsText = "after 1 wrong submission"
	} else {
		attemptsText = fmt.Sprintf("after %d wrong submissions", wrongAttempts)
	}

	// Convert Unix timestamp to UTC+2
	utcPlus2 := time.FixedZone("UTC+2", 2*60*60)
	solvedTime := time.Unix(sub.CreationTimeSeconds, 0).In(utcPlus2)
	timeStr := solvedTime.Format("2006-01-02 15:04:05")

	// Format task index with closing parenthesis if it looks like it should have one (e.g., "I" -> "I)")
	taskDisplay := sub.Problem.Index
	if len(sub.Problem.Index) == 1 && sub.Problem.Index[0] >= 'A' && sub.Problem.Index[0] <= 'Z' {
		taskDisplay = sub.Problem.Index + ")"
	}

	message := fmt.Sprintf(
		"🎉 <b>%s</b> (%s) has solved task <b>%s</b> <a href=\"%s\">%s</a> %s\n\n"+
			"🆗 Tests passed: %d\n"+
			"⏱ Time: %d ms\n"+
			"💾 Memory: %.2f MB\n"+
			"💻 Language: %s\n"+
			"🕐 Solved at: %s (UTC+2)\n"+
			"🔗 <a href=\"%s\">Submission #%d</a>",
		handle,
		realName,
		taskDisplay,
		problemURL,
		sub.Problem.Name,
		attemptsText,
		sub.PassedTestCount,
		sub.TimeConsumedMillis,
		float64(sub.MemoryConsumedBytes)/(1024*1024),
		sub.ProgrammingLanguage,
		timeStr,
		submissionURL,
		sub.ID,
	)

	return message
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	cfClient := NewCodeforcesClient()
	tracker := NewTracker()

	// Initial population of tracker to avoid spam on first run
	log.Println("Initializing tracker with recent submissions...")
	for _, h := range config.CodeforcesHandles {
		subs, err := cfClient.GetUserSubmissions(h.Handle, 20)
		if err != nil {
			log.Printf("Warning: failed to fetch initial submissions for %s: %v", h.Handle, err)
			continue
		}
		for _, sub := range subs {
			tracker.MarkProcessed(h.Handle, sub.ID)
		}
	}

	log.Println("Starting monitoring loop...")
	ticker := time.NewTicker(time.Duration(config.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, h := range config.CodeforcesHandles {
			subs, err := cfClient.GetUserSubmissions(h.Handle, 20)
			if err != nil {
				log.Printf("Error fetching submissions for %s: %v", h.Handle, err)
				continue
			}

			// Process submissions in reverse order (oldest first)
			for i := len(subs) - 1; i >= 0; i-- {
				sub := subs[i]

				if tracker.IsProcessed(h.Handle, sub.ID) {
					continue
				}

				tracker.MarkProcessed(h.Handle, sub.ID)

				// Only notify for accepted solutions
				if sub.Verdict != "OK" {
					continue
				}

				// Count wrong attempts for this problem
				problemKey := fmt.Sprintf("%d%s", sub.Problem.ContestID, sub.Problem.Index)
				wrongAttempts := tracker.CountWrongAttempts(h.Handle, problemKey, subs)

				message := formatNotification(h.Handle, h.RealName, sub, wrongAttempts)

				msg := tgbotapi.NewMessage(config.TelegramChatID, message)
				msg.ParseMode = "HTML"
				msg.DisableWebPagePreview = false

				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending message: %v", err)
				} else {
					log.Printf("Notified: %s solved %s%s", h.Handle, sub.Problem.Index, sub.Problem.Name)
				}
			}
		}
	}
}
