package structs

import (
	"encoding/json"
)

// Reputations contains the reputation for all available campaigns
type Reputations map[string]*Reputation

// Reputation describes the reputation progress for a specific campaign
type Reputation struct {
	Motto    string
	RankName *string `json:"Rank"`

	Level    *int
	Progress *float64

	ItemsTotal         int
	ItemsUnlocked      int
	EmblemsTotal       int
	EmblemsUnlocked    int
	TitlesTotal        *int
	TitlesUnlocked     *int
	PromotionsTotal    *int
	PromotionsUnlocked *int

	Campaigns map[string]Campaign
	Emblems   Emblems

	unlockSummaries map[string]UnlockSummary `json:"-"`
}

// UnlockSummary provides details about unlocks for a specific type (e.g. "items" or "emblems")
type UnlockSummary struct {
	Total    int
	Unlocked int
}

// UnlockSummaries maps each unlock type with its data.
// The returned map is cached in the struct.
func (r *Reputation) UnlockSummaries() map[string]UnlockSummary {
	if r.unlockSummaries == nil {
		s := make(map[string]UnlockSummary)
		// Items & Emblems are always present
		s["Items"] = UnlockSummary{Total: r.ItemsTotal, Unlocked: r.ItemsUnlocked}
		s["Emblems"] = UnlockSummary{Total: r.EmblemsTotal, Unlocked: r.EmblemsUnlocked}
		// Others optional
		if r.TitlesTotal != nil {
			s["Titles"] = UnlockSummary{Total: *r.TitlesTotal, Unlocked: *r.TitlesUnlocked}
		}
		if r.PromotionsTotal != nil {
			s["Promotions"] = UnlockSummary{Total: *r.PromotionsTotal, Unlocked: *r.PromotionsUnlocked}
		}
		r.unlockSummaries = s
	}
	return r.unlockSummaries
}

type Emblems []Emblem

// UnmarshalJSON unwraps the inner emblem list to avoid redundant nesting
func (l *Emblems) UnmarshalJSON(data []byte) (err error) {
	auxInnerData := struct {
		Emblems []Emblem
	}{}
	if err = json.Unmarshal(data, &auxInnerData); err != nil {
		return
	}
	*l = (Emblems)(auxInnerData.Emblems)
	return
}

type Emblem struct {
	Title           string
	Subtitle        string
	Locked          bool
	HasProgress     bool   `json:"HasScalar"` // true if Value/Threshold fields are applicable
	ProgressCurrent int    `json:"Value"`
	ProgressMax     int    `json:"Threshold"`
	ImageURL        string `json:"image"`
}

type Campaign struct {
	Title           string
	Subtitle        string `json:"Desc"`
	EmblemsTotal    int
	EmblemsUnlocked int
	Emblems         []Emblem
}
