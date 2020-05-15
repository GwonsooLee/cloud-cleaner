package reporter

import (
	"fmt"
	"github.com/slack-go/slack"
	"os"
	"strings"
)

var (
	START_TITLE="Hey there 👋 I'm *Clean Bot*..\nI'm here to help you check and clean wasted resources in Slack."
	GREAT_MESSAGE=":thumbsup::100: *No wasted resources* "
)

var (
	RegionMap=map[string]string{
		"us-east-1" : "N. Virginia",
		"us-west-2" : "Oregon",
		"ap-northeast-2" : "Seoul",
		"ap-southeast-1" : "Singapore",
		"us-east-2" : "Ohio",
		"af-south-1" : "Cape Town",
		"ap-east-1" : "Hong Kong",
		"ap-south-1" : "Mumbai",
		"ap-southeast-2" : "Sydney",
		"ap-northeast-1" : "Tokyo",
		"ca-central-1" : "Canada Central",
		"eu-central-1" : "Frankfurt",
		"eu-west-1" : "Ireland",
		"eu-west-2" : "London",
		"eu-south-1" : "Milan",
		"eu-west-3" : "Paris",
		"eu-north-1" : "Stockholm",
		"me-south-1" : "Bahrain",
		"sa-east-1" : "South America",
	}
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

func (r Reporter) SendGreatMessage() {
	textSection := r.CreateSimpleSection(GREAT_MESSAGE)
	msgOpt := slack.MsgOptionBlocks(textSection)
	r.SendMessage(msgOpt)
}


func (r Reporter) SendTitleMessage() {
	textSection := r.CreateSimpleSection(START_TITLE)
	msgOpt := slack.MsgOptionBlocks(textSection)
	r.SendMessage(msgOpt)
}

func (r Reporter) SendRegionMessage(region string) {
	msgOpt := r.CreateSimpleAttachments("Region", fmt.Sprintf("%s, %s", RegionMap[region], region))
	r.SendMessage(msgOpt)
}

func (r Reporter) CreateAlarmMessage(title string, sl []string) slack.MsgOption {
	return slack.MsgOptionBlocks(
		r.CreateDividerSection(),
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
	txt := slack.NewTextBlockObject("mrkdwn", text, false, false)
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

func (r Reporter) CreateSimpleAttachments(title, text string) slack.MsgOption {
	return slack.MsgOptionAttachments(
		slack.Attachment{
			Color:         "#36a64f",
			Title:         title,
			Text:          text,
			MarkdownIn:    []string{"text"},
		},
	)
}
