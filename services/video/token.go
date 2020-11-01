package video

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/pborman/uuid"
)

type AccessTokenHeader struct {
	Type        string `json:"typ"`
	Algorithm   string `json:"alg"`
	ContentType string `json:"cty"`
}

type VideoGrant struct {
	Room string `json:"room"`
}

type AccessTokenGrant struct {
	Identity string     `json:"identity"`
	Video    VideoGrant `json:"video"`
}

type AccessTokenBody struct {
	JTI    string           `json:"jti"`
	ISS    string           `json:"iss"`
	Sub    string           `json:"sub"`
	Nbf    int64            `json:"nbf"`
	Exp    int64            `json:"exp"`
	Grants AccessTokenGrant `json:"grants"`
}

// Generate the Rooms access token
func (c *VideoClient) GenerateAccessToken(uid string, room string) (string, error) {
	header := AccessTokenHeader{
		Type:        "JWT",
		Algorithm:   "HS256",
		ContentType: "twilio-fpa;v=1",
	}

	headerString, err := json.Marshal(&header)
	if err != nil {
		return "", err
	}

	videoGrant := VideoGrant{
		Room: room,
	}

	grant := AccessTokenGrant{
		Identity: uid,
		Video:    videoGrant,
	}

	nbf := time.Now()
	exp := nbf.Add(c.tokenExpiry)

	body := AccessTokenBody{
		JTI:    uuid.New(),
		ISS:    c.apiKey.SID,
		Sub:    c.accountSID,
		Nbf:    nbf.Unix(),
		Exp:    exp.Unix(),
		Grants: grant,
	}

	bodyString, err := json.Marshal(&body)
	if err != nil {
		return "", err
	}

	payload := base64.StdEncoding.EncodeToString(headerString) + "." + base64.StdEncoding.EncodeToString(bodyString)
	sig := hmac.New(sha256.New, []byte(c.apiKey.Secret))
	sig.Write([]byte(payload))

	return payload + "." + string(sig.Sum(nil)), nil
}
