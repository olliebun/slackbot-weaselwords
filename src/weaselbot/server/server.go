package server

import (
	"fmt"
	"os"
	"strings"

	"weaselbot"
	"weaselbot/config"
	"weaselbot/slack"
)

// Weaselbot server. Sits on an open websocket reading incoming webhooks from Slack.
type Server interface {
	Run() error
}

func NewServer(cfg config.Config, words weaselbot.Words, users weaselbot.Users) Server {
	return &server{cfg, words, users, nil, make(chan message, cfg.Message_Queue_Length)}
}

type errInvalidMessage struct {
	msg string
}

func (e errInvalidMessage) Error() string {
	return fmt.Sprintf("errInvalidMessage: %s", e.msg)
}

type message struct {
	User_Name string
	Channel   string
	Text      string
}

func (m message) String() string {
	return fmt.Sprintf("@%s in #%s: %q", m.User_Name, m.Channel, m.Text)
}

func (m message) validate() error {
	if m.User_Name == "" {
		return errInvalidMessage{"user_name not set"}
	} else if m.Channel == "" {
		return errInvalidMessage{"channel_name not set"}
	} else if m.Text == "" {
		return errInvalidMessage{"message text not set"}
	}
	return nil
}

type server struct {
	cfg      config.Config
	words    weaselbot.Words
	users    weaselbot.Users
	slack    slack.Slack
	incoming chan message
}

func (s *server) Run() error {
	// Start a slack real-time messaging session
	slack, err := slack.New(s.cfg.Slack_Token)
	if err != nil {
		return err
	}
	s.slack = slack

	for {
		var slack_msg map[string]interface{}
		if err := s.slack.GetMessage(&slack_msg); err != nil {
			fmt.Fprintf(os.Stderr, "server: failed to parse slack message: %s\n", err)
			continue
		}

		s.handle_message(slack_msg)
	}
}

func (s *server) handle_message(msg map[string]interface{}) {

	if msg["type"] != "message" {
		fmt.Printf("Got misc message from slack: %#v\n", msg)
		return
	}

	user := msg["user"].(string)
	text := msg["text"].(string)
	channel := msg["channel"].(string)

	if !s.users.Matches(user) {
		fmt.Printf("ignoring message for user %q\n", user)
		return
	}

	found := s.words.Matches(text)
	if len(found) == 0 {
		fmt.Printf("no word matches in message in chan %q for user %q, skipping\n", channel, user)
		return
	}

	channel_name, err := s.slack.GetChannelName(channel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get channel name for channel %q: %s", channel, err)
		return
	}

	restext := fmt.Sprintf("You used some weasel words in a message to #%s:\n```%v\n```", channel_name, strings.Join(found, "\n"))

	err = s.slack.SendDirectMessage(slack.DirectMessage{User_Name: user, Text: restext})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send slack DM: %s\n", err)
	}
}
