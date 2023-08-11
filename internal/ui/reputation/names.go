package reputation

import "image/color"

const (
	repTypeAthenasFortune   = "AthenasFortune"
	repTypeHuntersCall      = "HuntersCall"
	repTypeGoldHoarders     = "GoldHoarders"
	repTypeOrderOfSouls     = "OrderOfSouls"
	repTypeMerchantAlliance = "MerchantAlliance"
	repTypeCreatorCrew      = "CreatorCrew"
	repTypeBilgeRats        = "BilgeRats"
	repTypeTallTales        = "TallTales"
	repTypeReapersBones     = "ReapersBones"
	repTypePirateLord       = "PirateLord"
	repTypeFlameheart       = "Flameheart"
)

var repTypeDisplaynames = map[string]string{
	repTypeAthenasFortune:   "Athenas Fortune",
	repTypeHuntersCall:      "Hunter's Call",
	repTypeGoldHoarders:     "Gold Hoarders",
	repTypeOrderOfSouls:     "Order of Souls",
	repTypeMerchantAlliance: "Merchant Alliance",
	repTypeCreatorCrew:      "Creator Crew",
	repTypeBilgeRats:        "Bilge Rats",
	repTypeTallTales:        "Tall Tales",
	repTypeReapersBones:     "Reaper's Bones",
	repTypePirateLord:       "Pirate Lord",
	repTypeFlameheart:       "Flameheart",
}

const colorBgAlpha = 80

var repTypeColors = map[string]color.Color{
	repTypeAthenasFortune:   color.NRGBA{R: 30, G: 81, B: 98, A: colorBgAlpha},
	repTypeHuntersCall:      color.NRGBA{R: 80, G: 96, B: 99, A: colorBgAlpha},
	repTypeGoldHoarders:     color.NRGBA{R: 171, G: 154, B: 45, A: colorBgAlpha},
	repTypeOrderOfSouls:     color.NRGBA{R: 102, G: 46, B: 55, A: colorBgAlpha},
	repTypeMerchantAlliance: color.NRGBA{R: 19, G: 149, B: 178, A: colorBgAlpha},
	repTypeCreatorCrew:      color.NRGBA{R: 70, G: 70, B: 70, A: colorBgAlpha},
	repTypeBilgeRats:        color.NRGBA{R: 49, G: 9, B: 8, A: colorBgAlpha},
	repTypeTallTales:        color.NRGBA{R: 63, G: 64, B: 21, A: colorBgAlpha},
	repTypeReapersBones:     color.NRGBA{R: 202, G: 61, B: 48, A: colorBgAlpha},
	repTypePirateLord:       color.NRGBA{R: 30, G: 81, B: 98, A: colorBgAlpha},
	repTypeFlameheart:       color.NRGBA{R: 53, G: 21, B: 14, A: colorBgAlpha},
}
