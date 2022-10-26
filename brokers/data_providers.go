package brokers

import (
	"time"

	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/dto/user"
)

type DataProvider interface {
	AuthenticationUrl() string
	SetAccessToken(accessToken string) error
	Authenticate(authToken string) (*user.User, error)
	DownloadInstruments() (*instruments.Instruments, error)
	DownloadHistoricalData(ins instruments.Instruments, from, to time.Time, interval string) (*instruments.Instruments, error)
}
