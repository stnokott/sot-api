// Package api allows interaction with the SoT API
package api

import (
	"golang.org/x/exp/slog"
	"golang.org/x/text/language"

	"github.com/stnokott/sot-api/internal/api/structs"
)

// Client allows querying SoT API endpoints
type Client struct {
	token       string
	baseURL     string
	httpReferer string

	logger *slog.Logger
}

// NewClient creates a new API client with the specified token
func NewClient(ratToken string, locale language.Tag, logger *slog.Logger) *Client {
	var localeURLPart string
	if locale != language.English {
		localeURLPart = locale.String() + "/"
	}
	return &Client{
		token:       ratToken,
		baseURL:     "https://www.seaofthieves.com/" + localeURLPart + "api",
		httpReferer: "https://www.seaofthieves.com/" + localeURLPart + "profile",
		logger:      logger,
	}
}

// GetHealth retrieves data about the status of the API
func (c *Client) GetHealth() (h structs.Health, err error) {
	c.logger.Info("getting API health")
	err = c.get("https://status.seaofthieves.com/api/health", &h)
	return
}

// GetProfile retrieves the balance of in-game currencies plus title and profile image for the pirate
func (c *Client) GetProfile() (p structs.Profile, err error) {
	c.logger.Info("getting pirate profile")
	err = c.apiGet("/profilev2/balance", &p)
	return
}

// GetReputation retrieves the reputation for all available campaigns
func (c *Client) GetReputation() (r structs.Reputation, err error) {
	c.logger.Info("getting pirate reputation")
	err = c.apiGet("/profilev2/reputation", &r)
	return
}

// GetSeason retrieves data about the current active season
func (c *Client) GetSeason() (s structs.Season, err error) {
	c.logger.Info("getting pirate season progress")
	err = c.apiGet("/profilev2/seasons-progress", &s)
	return
}
