package structs

import (
	"encoding/json"
)

// Reputations contains the reputation for all available campaigns
type Reputations map[string]Reputation

type Reputation struct {
	Motto    string
	RankName *string `json:"Rank"`

	Level              *int
	Progress           *float64
	EmblemsTotal       *int
	EmblemsUnlocked    *int
	TitlesTotal        *int
	TitlesUnlocked     *int
	ItemsTotal         *int
	ItemsUnlocked      *int
	PromotionsTotal    *int
	PromotionsUnlocked *int

	Campaigns map[string]repCampaign
	Emblems   repEmblems
}

type repEmblems []repEmblem

// UnmarshalJSON unwraps the inner emblem list to avoid redundant nesting
func (l *repEmblems) UnmarshalJSON(data []byte) (err error) {
	auxInnerData := struct {
		Emblems []repEmblem
	}{}
	if err = json.Unmarshal(data, &auxInnerData); err != nil {
		return
	}
	*l = (repEmblems)(auxInnerData.Emblems)
	return
}

type repEmblem struct {
	Title           string
	Subtitle        string
	Locked          bool
	HasProgress     bool   `json:"HasScalar"` // true if Value/Threshold fields are applicable
	ProgressCurrent int    `json:"Value"`
	ProgressMax     int    `json:"Threshold"`
	ImageURL        string `json:"image"`
}

type repCampaign struct {
	Title           string
	Description     string `json:"Desc"`
	EmblemsTotal    int
	EmblemsUnlocked int
	Emblems         []repEmblem
}
