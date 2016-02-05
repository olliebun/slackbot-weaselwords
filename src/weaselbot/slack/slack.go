package slack

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

const (
	DM_CHANNEL_ID = "D024BE91L"
	SLACK_API     = "https://api.slack.com"
)

type DirectMessage struct {
	User_Name string
	Text      string
}

type slackDirectMessage struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type errSlackSetupFailed struct {
	msg string
}

func (e errSlackSetupFailed) Error() string {
	return fmt.Sprintf("errSlackSetupFailed: %s", e.msg)
}

type Slack interface {
	GetMessage(interface{}) error
	SendDirectMessage(DirectMessage) error

	GetChannelName(channel_id string) (string, error)
	GetUserName(user_id string) (string, error)
}

type slack struct {
	conn *websocket.Conn

	slack_token string

	// Map of user ID to IM channel IDs
	imChannels map[string]string

	// Map of channel ID to channel name
	channels map[string]string

	// Map of user ID to user name
	users map[string]string
}

func (s *slack) slack_get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	var res map[string]interface{}
	url := fmt.Sprintf("%s/api/%s?token=%s", SLACK_API, endpoint, s.slack_token)

	for k, v := range params {
		url = url + "&" + k + "=" + v
	}

	fmt.Printf("Hitting slack API at %s\n", url)
	resp, err := http.Get(url)

	if err != nil {
		return nil, errSlackSetupFailed{err.Error()}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errSlackSetupFailed{fmt.Sprintf("Response status code from slack API was %d", resp.StatusCode)}
	}

	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(&res); err != nil {
		return nil, errSlackSetupFailed{fmt.Sprintf("Failed to decode response from slack API: %s", err)}
	}

	fmt.Printf("Response from slack API: %#v\n", res)

	if status, ok := res["ok"]; ok && status.(bool) == false {
		return nil, errSlackSetupFailed{fmt.Sprintf(`Got response with "ok": false from %s`, endpoint)}

	}

	return res, nil
}

// Open an IM channel with this user. Returns the channel ID.
func (s *slack) im_open(user_id string) (string, error) {
	res, err := s.slack_get("im.open", map[string]string{"user": user_id})
	fmt.Printf("Response from slack: %#v\n", res)

	if err != nil {
		return "", err
	}

	channel_id := res["channel"].(map[string]interface{})["id"].(string)
	s.imChannels[user_id] = channel_id
	return channel_id, nil
}

func (s *slack) update_channel_list() error {
	res, err := s.slack_get("channels.list", nil)

	if err != nil {
		return err
	}

	for _, c := range res["channels"].([]interface{}) {
		cmap := c.(map[string]interface{})
		s.channels[cmap["id"].(string)] = cmap["name"].(string)
	}

	return nil
}

func (s *slack) update_user_list() error {
	res, err := s.slack_get("users.list", nil)

	if err != nil {
		return err
	} else if !res["ok"].(bool) {
		return errSlackSetupFailed{"Failed to get channel list"}
	}

	for _, c := range res["members"].([]interface{}) {
		cmap := c.(map[string]interface{})
		s.users[cmap["id"].(string)] = cmap["name"].(string)
	}

	return nil
}

// Given a slack token, authenticate with to the RTM API and return the websocket URL.
func slack_api_connect(slack_token string) (string, error) {
	url := fmt.Sprintf("%s/api/rtm.start?token=%s", SLACK_API, slack_token)
	resp, err := http.Get(url)
	if err != nil {
		return "", errSlackSetupFailed{err.Error()}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errSlackSetupFailed{fmt.Sprintf("Response status code from slack API was %d", resp.StatusCode)}
	}

	dec := json.NewDecoder(resp.Body)
	var res map[string]interface{}

	if err := dec.Decode(&res); err != nil {
		return "", errSlackSetupFailed{fmt.Sprintf("Failed to decode response from slack API: %s", err)}
	}

	if es, ok := res["error"]; ok {
		return "", errSlackSetupFailed{es.(string)}
	}

	url, ok := res["url"].(string)
	if !ok {
		return "", errSlackSetupFailed{fmt.Sprintf("Expected slack websocket URL to be a string; got %T", res["url"])}
	}

	return url, nil
}

func slack_websocket_connect(url string) (*websocket.Conn, error) {
	return websocket.Dial(url, "", SLACK_API)
}

func New(slack_token string) (Slack, error) {
	url, err := slack_api_connect(slack_token)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully authenticated with slack; connecting to websocket url: %s\n", url)
	conn, err := slack_websocket_connect(url)
	if err != nil {
		return nil, err
	}

	s := &slack{conn, slack_token, make(map[string]string), make(map[string]string), make(map[string]string)}

	return s, nil
}

func (s *slack) GetMessage(into interface{}) error {
	err := websocket.JSON.Receive(s.conn, into)
	return err
}

func (s *slack) SendDirectMessage(msg DirectMessage) (err error) {
	// Get the channel id
	channel_id, ok := s.imChannels[msg.User_Name]
	if !ok {
		fmt.Printf("No channel ID found for user %s: looking one up\n", msg.User_Name)
		channel_id, err = s.im_open(msg.User_Name)
		if err != nil {
			return err
		}

	}

	slackmsg := slackDirectMessage{Id: 1, Type: "message", Channel: channel_id, Text: msg.Text}
	// slackmsg := slackDirectMessage{Id: 1, Type: "message", Channel: "#general", Text: msg.Text}

	fmt.Printf("Sending message to slack: %#v\n", slackmsg)
	err = websocket.JSON.Send(s.conn, slackmsg)
	if err != nil {
		return err
	}

	// TODO(cera) - Read the response from Slack here

	return nil
}

func (s *slack) GetChannelName(channel_id string) (string, error) {
	channel_name, ok := s.channels[channel_id]
	if ok {
		return channel_name, nil
	}

	err := s.update_channel_list()
	if err != nil {
		return "", err
	}

	channel_name, ok = s.channels[channel_id]
	if !ok {
		return "", fmt.Errorf("Channel name not found")
	}
	return channel_name, nil
}

func (s *slack) GetUserName(user_id string) (string, error) {
	user_name, ok := s.users[user_id]
	if ok {
		return user_name, nil
	}

	err := s.update_user_list()
	if err != nil {
		return "", err
	}

	user_name, ok = s.users[user_id]
	if !ok {
		return "", fmt.Errorf("User name not found")
	}
	return user_name, nil
}
