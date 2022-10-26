package zerodha

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tekramkcots/sdk/dto/user"
	"github.com/tekramkcots/sdk/markets/indian"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

type Client struct {
	*kiteconnect.Client
	apiSecret string
	logger    *logrus.Entry
}

func New(apiKey, apiSecret string, logger *logrus.Entry) (*Client, error) {
	kc := kiteconnect.New(apiKey)
	return &Client{kc, apiSecret, logger}, nil
}

func (c *Client) AuthenticationUrl() string {
	return c.GetLoginURL()
}

func (c *Client) SetAccessToken(accessToken string) error {
	c.Client.SetAccessToken(accessToken)
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, indian.GetTimeZone()).AddDate(0, 0, 1)
	go time.AfterFunc(nextMidnight.Sub(now), func() {
		c.SetAccessToken("")
	})
	return nil
}

func (c *Client) Authenticate(authToken string) (*user.User, error) {
	data, err := c.Client.GenerateSession(authToken, c.apiSecret)
	if err != nil {
		return nil, fmt.Errorf("error generating the session %w", err)
	}
	c.SetAccessToken(data.AccessToken)
	profile, err := c.Client.GetUserProfile()
	if err != nil {
		return nil, fmt.Errorf("error getting user profile %w", err)
	}
	return &user.User{
		UserID:      profile.UserID,
		UserName:    profile.UserName,
		AccessToken: data.AccessToken,
		Email:       profile.Email,
	}, nil
}
