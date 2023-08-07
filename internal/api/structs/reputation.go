package structs

import (
	"encoding/json"
)

// Reputations contains the reputation for all available campaigns
type Reputations map[string]Reputation

type Reputation struct {
	Motto    string
	RankName *string `json:"Rank"`

	Level    *int
	Progress *float64

	UnlockSummaries map[string]unlockSummary `json:"-"`

	Campaigns map[string]repCampaign
	Emblems   repEmblems

	unlocks json.RawMessage
}

func (r *Reputation) UnmarshalJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, r); err != nil {
		return
	}

	aux := struct {
		ItemsTotal         int
		ItemsUnlocked      int
		EmblemsTotal       int
		EmblemsUnlocked    int
		TitlesTotal        *int
		TitlesUnlocked     *int
		PromotionsTotal    *int
		PromotionsUnlocked *int
	}{}
	if err = json.Unmarshal(r.unlocks, &aux); err != nil {
		return
	}
	// Items & Emblems are always present
	r.UnlockSummaries["Items"] = unlockSummary{Total: aux.ItemsTotal, Unlocked: aux.ItemsUnlocked}
	r.UnlockSummaries["Emblems"] = unlockSummary{Total: aux.EmblemsTotal, Unlocked: aux.EmblemsUnlocked}
	// Others optional
	if aux.TitlesTotal != nil {
		r.UnlockSummaries["Titles"] = unlockSummary{Total: *aux.TitlesTotal, Unlocked: *aux.TitlesUnlocked}
	}
	if aux.PromotionsTotal != nil {
		r.UnlockSummaries["Promotions"] = unlockSummary{Total: *aux.PromotionsTotal, Unlocked: *aux.PromotionsUnlocked}
	}
	return
}

type unlockSummary struct {
	Total    int
	Unlocked int
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
