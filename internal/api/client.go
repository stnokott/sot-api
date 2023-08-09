// Package api allows interaction with the SoT API
package api

import (
	"go.uber.org/zap"
	"golang.org/x/text/language"

	"github.com/stnokott/sot-api/internal/api/structs"
	"github.com/stnokott/sot-api/internal/files"
)

// Client allows querying SoT API endpoints
type Client struct {
	token       string
	baseURL     string
	httpReferer string

	logger *zap.Logger
}

// NewClient creates a new API client with the specified token
func NewClient(token string, locale language.Tag, logger *zap.Logger) *Client {
	var localeURLPart string
	if locale != language.English {
		// URL for English does not contain a locale string, all others do
		localeURLPart = locale.String() + "/"
	}

	return &Client{
		token:       token,
		baseURL:     "https://www.seaofthieves.com/" + localeURLPart + "api",
		httpReferer: "https://www.seaofthieves.com/" + localeURLPart + "profile",
		logger:      logger,
	}
}

// SetToken updates the RAT token and writes it to file
func (c *Client) SetToken(t string) {
	c.token = t
	if err := files.WriteToken(t); err != nil {
		c.logger.Fatal("could not write token: " + err.Error())
	}
}

// GetHealth retrieves data about the status of the API
func (c *Client) GetHealth() (h *structs.Health, err error) {
	c.logger.Info("getting API health")
	err = c.get("https://status.seaofthieves.com/api/health", &h)
	return
}

// GetProfile retrieves the balance of in-game currencies plus title and profile image for the pirate
func (c *Client) GetProfile() (p *structs.Profile, err error) {
	c.logger.Info("getting profile")
	p = new(structs.Profile)
	err = c.apiGet("/profilev2/balance", p)
	return
}

// GetReputation retrieves the reputation for all available campaigns
func (c *Client) GetReputation() (r structs.Reputations, err error) {
	c.logger.Info("getting reputation")
	err = c.apiGet("/profilev2/reputation", &r)
	return
}

// GetSeason retrieves data about the current active season
func (c *Client) GetSeason() (s *structs.Season, err error) {
	c.logger.Info("getting season progress")
	s = new(structs.Season)
	err = c.apiGet("/profilev2/seasons-progress", s)
	return
}
