package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/msg"

	"github.com/mattermost/mattermost-server/v5/model"
)

// The Mattermost exporter sends Flux events to a Mattermost channel via a webhook.
type Mattermost struct {
	Url      string
	Username string
	Channels []MattermostChannel
	IconURL  string
}

// MattermostChannel Represents a mattermosts channel and the Kubernetes namespace linked to it
type MattermostChannel struct {
	Channel   string `json:"channel"`
	Namespace string `json:"namespace"`
}

// NewMattermost initialize a new Mattermost instance
func NewMattermost(config config.Config) (*Mattermost, error) {
	var err error
	m := Mattermost{}

	m.Url, err = config.Required("mattermost_url")
	if err != nil {
		return nil, err
	}

	channels, err := config.Required("mattermost_channel")
	if err != nil {
		return nil, err
	}
	m.parseMattermostChannelConfig(channels)
	log.Println(m.Channels)

	m.Username = config.Optional("mattermost_username", "Flux Deployer")
	m.IconURL = config.Optional("mattermost_icon_url", "https://user-images.githubusercontent.com/27962005/35868977-0d5f85f6-0b2c-11e8-9fa8-8e4eaf35161a.png")

	return &m, nil
}

// Send a MattermostMessage to Mattermost
func (m *Mattermost) Send(c context.Context, client *http.Client, message msg.Message) error {
	for _, mattermostMessage := range m.NewMattermostMessage(message) {
		fmt.Println(mattermostMessage)
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(mattermostMessage)
		if err != nil {
			log.Print("Could encode message to mattermost:", err)
			return err
		}

		log.Print(string(b.Bytes()))

		req, _ := http.NewRequest("POST", m.Url, b)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(c)

		res, err := client.Do(req)
		if err != nil {
			log.Print("Could not post to mattermost:", err)
			return err
		}

		if res.StatusCode != 200 {
			log.Print("Could not post to mattermost, status: ", res.Status)
			return fmt.Errorf("Could not post to mattermost, status: %d", res.StatusCode)
		}
	}

	return nil
}

// NewLine return the new line character for Slack messages
func (m *Mattermost) NewLine() string {
	return "\n"
}

// FormatLink return a formatted link for Slack.
func (m *Mattermost) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

// NewMattermostMessage convert a flux event into a Mattermost message(s)
func (m *Mattermost) NewMattermostMessage(message msg.Message) []model.IncomingWebhookRequest {
	var messages []model.IncomingWebhookRequest

	for _, channel := range m.determineChannels(message) {
		mattermostMessage := model.IncomingWebhookRequest{
			ChannelName: channel,
			IconURL:     m.IconURL,
			Username:    m.Username,
			Attachments: []*model.SlackAttachment{
				{
					Color:     "#4286f4",
					TitleLink: message.TitleLink,
					Title:     message.Title,
					Text:      message.Body,
				},
			},
		}
		messages = append(messages, mattermostMessage)
	}

	return messages
}

// Name return the name of the exporter.
func (m *Mattermost) Name() string {
	return "Mattermost"
}

// Parse the channel configuration string in a backwards
// compatible manner.
func (m *Mattermost) parseMattermostChannelConfig(channels string) error {
	if len(strings.Split(channels, "=")) == 1 {
		m.Channels = append(m.Channels, MattermostChannel{channels, "*"})
		return nil
	}

	re := regexp.MustCompile("([#a-z0-9][a-z0-9._-]*)=([a-z0-9*][-A-Za-z0-9_.]*)")
	for _, kv := range strings.Split(channels, ",") {
		if !re.MatchString(kv) {
			return fmt.Errorf("Could not parse channel/namespace configuration: %s", kv)
		}

		cn := strings.Split(kv, "=")
		channel := strings.TrimSpace(cn[0])
		namespace := strings.TrimSpace(cn[1])
		m.Channels = append(m.Channels, MattermostChannel{channel, namespace})
	}

	return nil
}

// Match namespaces from service IDs to Mattermost channels.
func (m *Mattermost) determineChannels(message msg.Message) []string {
	var channels []string
	for _, serviceID := range message.Event.ServiceIDs {
		ns, _, _ := serviceID.Components()

		for _, ch := range m.Channels {
			if ch.Namespace == "*" || ch.Namespace == ns {
				channels = appendIfMissing(channels, ch.Channel)
			}
		}
	}
	return channels
}
