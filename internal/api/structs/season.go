//go:generate stringer -type=xpGain -linecomment

package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Season contains data about the current active season
type Season struct {
	ActiveFrom  time.Time
	ActiveUntil time.Time

	UnlockedTiers struct {
		Base       bool
		Legendary  bool
		SeasonPass bool
	} `json:"AvailablePaths"`

	TotalChallenges    int
	CompleteChallenges int
	ChallengeGroups    []seasonChallengeGroup

	Tier  int
	Tiers []seasonTier

	CurrentLevelProgress float64 `json:"LevelProgress"`

	CdnPath string // URL prefix for resources such as reward images
}

// UnmarshalJSON filters out the current season from the list of past seasons
func (s *Season) UnmarshalJSON(data []byte) (err error) {
	var auxInnerData []struct {
		Season
		IsActive bool
	}
	if err = json.Unmarshal(data, &auxInnerData); err != nil {
		return
	}
	for _, season := range auxInnerData {
		if season.IsActive {
			*s = season.Season
			return
		}
	}
	err = errors.New("no active season found")
	return
}

type seasonTier struct {
	Title  string
	Levels []struct {
		Rewards struct {
			Base       []seasonTierReward
			Legendary  []seasonTierReward
			SeasonPass []seasonTierReward
		} `json:"RewardsV2"`
	}
}

type seasonTierReward struct {
	Title       string `json:"EntitlementText"`
	Description string `json:"EntitlementDescription"`
	ImageName   string `json:"EntitlementURL"` // contains only the image name, needs to be prefixed with CdnPath from `Season`

	Locked bool
	Owned  bool
}

type seasonChallengeGroup struct {
	Title    string
	Subtitle string `json:"Copy"`

	Challenges []seasonChallenge

	ChallengesCompleted int `json:"ProgressValue"`
	IsCompleted         bool
}

type xpGain int

const (
	xpGainS xpGain = iota // S
	xpGainM               // M
	xpGainL               // L
)

var xpGainMap = map[string]xpGain{
	"s": xpGainS,
	"l": xpGainL,
	"m": xpGainM,
}

func (x *xpGain) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("XPGain should be a string, got %s", data)
	}
	val, ok := xpGainMap[s]
	if !ok {
		return fmt.Errorf("'%s' is not a valid XPGain size", s)
	}
	*x = val
	return nil
}

type seasonChallenge struct {
	Title    string
	Subtitle string `json:"Copy"`
	XPGain   xpGain `json:"XPGain"`

	Goals []seasonChallengeGoal

	GoalsCompleted int `json:"ProgressValue"`
	IsCompleted    bool
}

type seasonChallengeGoal struct {
	Title  string
	XPGain xpGain `json:"XPGain"`

	ProgressCurrent int `json:"ProgressValue"`
	ProgressMax     int `json:"Threshold"`
}
