package telegram

import (
	"fmt"
	"log"
	"time"

	"cfvgtuai/internal/codeforces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Notifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewNotifier(token string, chatID int64) (*Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Notifier{bot: bot, chatID: chatID}, nil
}

func (n *Notifier) Notify(submission codeforces.Submission, handle string, wrongAttempts int) {
	problemURL := ""
	if submission.Problem.Type == "PROGRAMMING" {
		problemURL = fmt.Sprintf("https://codeforces.com/contest/%d/problem/%s", submission.Problem.ContestID, submission.Problem.Index)
	} else {
		problemURL = fmt.Sprintf("https://codeforces.com/gym/%d/problem/%s", submission.Problem.ContestID, submission.Problem.Index)
	}

	submissionTime := time.Unix(submission.CreationTimeSeconds, 0).UTC().Add(2 * time.Hour)

	message := fmt.Sprintf(
		"User *%s* has solved problem [%s](%s)!\n\n"+
			"Wrong attempts before solving: %d\n"+
			"Tests: %d\n"+
			"Time: %d ms\n"+
			"Memory: %.2f MB\n"+
			"Language: %s\n"+
			"[Submission](https://codeforces.com/contest/%d/submission/%d) | %s",
		handle,
		submission.Problem.Name,
		problemURL,
		wrongAttempts,
		submission.PassedTestCount,
		submission.TimeConsumedMillis,
		float64(submission.MemoryConsumedBytes)/1024/1024,
		submission.ProgrammingLanguage,
		submission.ContestID,
		submission.ID,
		submissionTime.Format("2006-01-02 15:04:05"),
	)

	msg := tgbotapi.NewMessage(n.chatID, message)
	msg.ParseMode = "Markdown"

	if _, err := n.bot.Send(msg); err != nil {
		log.Printf("failed to send telegram notification: %v", err)
	}
}
