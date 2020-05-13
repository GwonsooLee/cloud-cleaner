package reporter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Reporter struct {
	WebhookUrl string
	Token string
}

type SlackBody struct {
	Text string `json:"text"`
}

func New(url string, token string) Reporter {
	return Reporter{
		WebhookUrl: url,
		Token:     token,
	}
}

func (r Reporter) Send_slack_message(message string) {
	slackBody, _ := json.Marshal(SlackBody{Text: message})
	req, err := http.NewRequest(http.MethodPost, r.WebhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		os.Exit(1)
	}
}
