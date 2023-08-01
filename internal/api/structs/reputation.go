package structs

import (
	"encoding/json"
)

// Reputation contains the reputation for all available campaigns
type Reputation struct {
	AthenasFortune   *repAthenasFortune
	HuntersCall      *repHuntersCall
	GoldHoarders     *repGoldHoarders
	OrderOfSouls     *repOrderOfSouls
	MerchantAlliance *repMerchantAlliance
	CreatorCrew      *repCreatorCrew
	BilgeRats        *repBilgeRats
	TallTales        *repTallTales
	ReapersBones     *repReapersBones
	PirateLord       *repPirateLord
	Flameheart       *repFlameheart
}

type repAthenasFortune struct {
	Emblems repEmblems

	repWithMotto
	repWithRank
	repWithLevel
	repWithEmblems
	repWithItems
	repWithTitles
}

func (r repAthenasFortune) String() string {
	return "AthenasFortune{Motto=" + r.Motto + "}"
}

type repHuntersCall struct {
	Campaigns map[string]repCampaign

	repWithMotto
	repWithRank
	repWithLevel
	repWithEmblems
	repWithItems
	repWithPromotions
	repWithTitles
}

func (r repHuntersCall) String() string {
	return "HuntersCall{Motto=" + r.Motto + "}"
}

type repGoldHoarders struct {
	Emblems repEmblems

	repWithMotto
	repWithRank
	repWithLevel
	repWithEmblems
	repWithItems
	repWithPromotions
	repWithTitles
}

func (r repGoldHoarders) String() string {
	return "GoldHoarders{Motto=" + r.Motto + "}"
}

type repOrderOfSouls struct {
	Emblems repEmblems

	repWithMotto
	repWithRank
	repWithLevel
	repWithEmblems
	repWithItems
	repWithPromotions
	repWithTitles
}

func (r repOrderOfSouls) String() string {
	return "OrderOfSouls{Motto=" + r.Motto + "}"
}

type repMerchantAlliance struct {
	Emblems repEmblems

	repWithMotto
	repWithRank
	repWithLevel
	repWithEmblems
	repWithItems
	repWithPromotions
	repWithTitles
}

func (r repMerchantAlliance) String() string {
	return "MerchantAlliance{Motto=" + r.Motto + "}"
}

type repCreatorCrew struct {
	Emblems repEmblems

	repWithMotto
	repWithEmblems
	repWithTitles
}

func (r repCreatorCrew) String() string {
	return "CreatorCrew{Motto=" + r.Motto + "}"
}

type repBilgeRats struct {
	Campaigns map[string]repCampaign

	repWithMotto
	repWithEmblems
	repWithTitles
}

func (r repBilgeRats) String() string {
	return "BilgeRats{Motto=" + r.Motto + "}"
}

type repTallTales struct {
	Campaigns map[string]repCampaign

	repWithMotto
	repWithTitles
	repWithEmblems
}

func (r repTallTales) String() string {
	return "TallTales{Motto=" + r.Motto + "}"
}

type repReapersBones struct {
	Emblems repEmblems

	repWithMotto
	repWithRank
	repWithLevel
	repWithPromotions
	repWithTitles
	repWithEmblems
	repWithItems
}

func (r repReapersBones) String() string {
	return "ReapersBones{Motto=" + r.Motto + "}"
}

type repPirateLord struct {
	Emblems repEmblems

	repWithMotto
	repWithLevel
	repWithTitles
	repWithEmblems
	repWithItems
}

func (r repPirateLord) String() string {
	return "PirateLord{Motto=" + r.Motto + "}"
}

type repFlameheart struct {
	Emblems repEmblems

	repWithMotto
	repWithLevel
	repWithTitles
	repWithEmblems
	repWithItems
}

func (r repFlameheart) String() string {
	return "Flameheart{Motto=" + r.Motto + "}"
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

type repWithMotto struct {
	Motto string
}

type repWithRank struct {
	RankName string `json:"Rank"`
}

type repWithEmblems struct {
	EmblemsTotal    int
	EmblemsUnlocked int
}

type repWithTitles struct {
	TitlesTotal    int
	TitlesUnlocked int
}

type repWithItems struct {
	ItemsTotal    int
	ItemsUnlocked int
}

type repWithPromotions struct {
	PromotionsTotal    int
	PromotionsUnlocked int
}

type repWithLevel struct {
	Level    int
	Progress float64
}
