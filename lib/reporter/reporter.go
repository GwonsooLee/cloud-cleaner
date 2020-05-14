package reporter

import (
	"os"
	"strings"
	"github.com/slack-go/slack"
)

type Reporter struct {
	Token string
	ChannelId string
}

type SlackBody struct {
	Text string `json:"text"`
}

func New(token string, channel_id string) Reporter {
	return Reporter{
		Token:     token,
		ChannelId:  channel_id,
	}
}

func (r Reporter) SendSimpleMessage(message string) {
	textSection := r.CreateSimpleSection(message)
	msgOpt := slack.MsgOptionBlocks(textSection)
	r.SendMessage(msgOpt)
}


func (r Reporter) CreateVolumeAlarmMessage(title string, sl []string) slack.MsgOption {
	return slack.MsgOptionBlocks(
		r.CreateTitleSection(title),
		r.CreateDividerSection(),
		r.CreateSimpleSection(strings.Join(sl[:], "\n")),
	)
}

func (r Reporter) SendBlockMessage(msgOpt slack.MsgOption) (string, string, string) {
	client := r.GetSlackClient()
	channel, timestamp, response, err:= client.SendMessage(r.ChannelId, msgOpt)
	if err != nil {
		os.Exit(1)
	}

	return channel, timestamp, response
}

func (r Reporter) SendMessage(msgOpt slack.MsgOption) (string, string, string) {
	client := r.GetSlackClient()
	channel, timestamp, response, err:= client.SendMessage(r.ChannelId, msgOpt)
	if err != nil {
		os.Exit(1)
	}

	return channel, timestamp, response
}

func (r Reporter) CreateSimpleSection(text string) *slack.SectionBlock {
	txt := slack.NewTextBlockObject("mrkdwn", text, false, false)
	section := slack.NewSectionBlock(txt, nil,nil)
	return section
}

func (r Reporter) CreateTitleSection(text string) *slack.SectionBlock {
	txt := slack.NewTextBlockObject("mrkdwn", "*"+text+"*", false, false)
	section := slack.NewSectionBlock(txt, nil,nil)
	return section
}

func (r Reporter) CreateDividerSection() *slack.DividerBlock {
	return slack.NewDividerBlock()
}

func (r Reporter) GetSlackClient() *slack.Client {
	token := r.Token
	if token == "" {
		token = os.Getenv("SLACK_OAUTH_TOKEN")
	}
	client := slack.New(token)
	return client
}
