package video

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const TWILIO_API_URL = "https://api.twilio.com/2010-04-01"
const TWILIO_VIDEO_API_URL = "https://video.twilio.com/v1"

type TwilioLink struct {
	Participants string `json:"participants"`
	Recordings   string `json:"recordings"`
}

type TwilioRoomResponse struct {
	AccountSID                  string       `json:"account_sid"`
	DateCreated                 time.Time    `json:"date_created"`
	DateUpdated                 time.Time    `json:"date_updated"`
	Status                      string       `json:"status"`
	Type                        string       `json:"type"`
	SID                         string       `json:"sid"`
	EnableTurn                  bool         `json:"enable_turn"`
	UniqueName                  string       `json:"unique_name"`
	MaxParticipants             int          `json:"max_participants"`
	Duration                    int          `json:"duration"`
	StatusCallbackMethod        string       `json:"status_callback_method"`
	StatusCallback              string       `json:"status_callback"`
	RecordParticipantsOnConnect bool         `json:"record_participants_on_connect"`
	VideoCodecs                 []string     `json:"video_codecs"`
	MediaRegion                 string       `json:"media_region"`
	EndTime                     time.Time    `json:"end_time"`
	Url                         string       `json:"url"`
	Links                       []TwilioLink `json:"links"`
}

type TwilioException struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"more_info"`
}

func (e *TwilioException) Error() string {
	return fmt.Sprintf("Status %d: Twilio error %d: %s; %s", e.Status, e.Code, e.Message, e.MoreInfo)
}

type APIKeyResponse struct {
	SID          string    `json:"sid"`
	FriendlyName string    `json:"friendly_name"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
	Secret       string    `json:"secret"`
}

type VideoClient struct {
	accountSID  string
	authToken   string
	client      *http.Client
	apiKey      *APIKeyResponse
	tokenExpiry time.Duration
}

// Initialise the video client. Requests an API key from the API.
func NewVideoClient(accountSID string, authToken string, tokenExpiry time.Duration) (*VideoClient, error) {
	client := &http.Client{}
	vc := &VideoClient{
		accountSID,
		authToken,
		client,
		nil,
		tokenExpiry,
	}

	key, err := vc.requestAPIKey()
	if err != nil {
		return nil, err
	}

	vc.apiKey = key

	return vc, nil
}

// Call the Rooms API to request an API key. Bootstrapping~
func (c *VideoClient) requestAPIKey() (*APIKeyResponse, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/Accounts/%s/Keys.json", TWILIO_API_URL, c.accountSID), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.accountSID, c.authToken)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusCreated {
		var body APIKeyResponse
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, err
		}

		return &body, nil
	} else {
		var body TwilioException
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, err
		}

		return nil, &body
	}
}

// Helper method to make HTTP requests to the Rooms API with the auth token in the Basic Auth. Also handles error codes.
func (c *VideoClient) makeRequest(method string, path string, values url.Values) (*TwilioRoomResponse, error) {
	var body string
	if values != nil {
		body = values.Encode()
	} else {
		body = ""
	}

	req, err := http.NewRequest(method, TWILIO_VIDEO_API_URL+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.accountSID, c.authToken)
	if values != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		var body TwilioRoomResponse
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, err
		}

		return &body, nil
	} else {
		var body TwilioException
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, err
		}

		return nil, &body
	}
}

// Calls the Rooms API to create a room.
func (c *VideoClient) CreateRoom(name string) (*TwilioRoomResponse, error) {
	v := url.Values{}
	if name != "" {
		v.Set("uniqueName", name)
	}

	res, err := c.makeRequest("POST", "/Rooms", v)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Calls the Rooms API to complete a room.
func (c *VideoClient) CompleteRoom(room string) error {
	values := url.Values{}
	values.Set("Status", "completed")
	_, err := c.makeRequest("POST", "/Rooms/"+room, values)
	if err != nil {
		return err
	}

	return nil
}

// Calls the Rooms API to get a room's information. Used for creating access tokens.
func (c *VideoClient) GetRoom(room string) (*TwilioRoomResponse, error) {
	res, err := c.makeRequest("GET", "/Rooms/"+room, nil)
	if err != nil {
		return nil, err
	}

	return res, err
}
