package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

type profileToolbar struct {
	title        *widget.Label
	gold         *widget.Label
	doubloons    *widget.Label
	ancientCoins *widget.Label

	widget.BaseWidget
}

func newProfileToolbar() *profileToolbar {
	title := widget.NewLabel("n/a")
	title.TextStyle.Bold = true
	title.TextStyle.Italic = true

	gold := widget.NewLabel("n/a")
	doubloons := widget.NewLabel("n/a")
	ancientCoins := widget.NewLabel("n/a")

	t := &profileToolbar{
		title:        title,
		gold:         gold,
		doubloons:    doubloons,
		ancientCoins: ancientCoins,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *profileToolbar) CreateRenderer() fyne.WidgetRenderer {
	iconMinSize := fyne.NewSize(16, 16)

	return widget.NewSimpleRenderer(
		container.NewHBox(
			t.title,
			newImageFromResource(resourceGoldPng, iconMinSize),
			t.gold,
			newImageFromResource(resourceDoubloonsPng, iconMinSize),
			t.doubloons,
			newImageFromResource(resourceAncientCoinsPng, iconMinSize),
			t.ancientCoins,
		),
	)
}

// SetProfile sets all labels and images according to the provided Profile
func (t *profileToolbar) SetProfile(p *structs.Profile) {
	t.title.SetText(p.Title)
	t.gold.SetText(strconv.Itoa(p.Gold))
	t.ancientCoins.SetText(strconv.Itoa(p.AncientCoins))
	t.doubloons.SetText(strconv.Itoa(p.Doubloons))
}
