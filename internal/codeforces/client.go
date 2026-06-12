package codeforces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Submission struct {
	ID                  int     `json:"id"`
	ContestID           int     `json:"contestId"`
	Problem             Problem `json:"problem"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict"`
	PassedTestCount     int     `json:"passedTestCount"`
	TimeConsumedMillis  int     `json:"timeConsumedMillis"`
	MemoryConsumedBytes int     `json:"memoryConsumedBytes"`
	CreationTimeSeconds int64   `json:"creationTimeSeconds"`
}

type Problem struct {
	ContestID int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type response struct {
	Status string       `json:"status"`
	Result []Submission `json:"result"`
}

type Client struct {
	lastRequest time.Time
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetUserSubmissions(handle string) ([]Submission, error) {
	since := time.Since(c.lastRequest)
	if since < 2*time.Second {
		time.Sleep(2*time.Second - since)
	}

	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=50", handle)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.lastRequest = time.Now()

	var apiResp response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("codeforces api error: %s", apiResp.Status)
	}

	return apiResp.Result, nil
}
