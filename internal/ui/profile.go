package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

type profileToolbar struct {
	title        binding.String
	gold         binding.Int
	doubloons    binding.Int
	ancientCoins binding.Int

	widget.BaseWidget
}

func (t *profileToolbar) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	lblTitle := widget.NewLabelWithData(t.title)
	lblTitle.TextStyle.Bold = true
	lblTitle.TextStyle.Italic = true

	iconMinSize := fyne.NewSize(16, 16)

	return widget.NewSimpleRenderer(
		container.NewHBox(
			lblTitle,
			widget.NewLabelWithData(binding.IntToString(t.ancientCoins)),
			newImageFromResource(resourceAncientCoinsPng, iconMinSize),
			widget.NewLabelWithData(binding.IntToString(t.doubloons)),
			newImageFromResource(resourceDoubloonsPng, iconMinSize),
			widget.NewLabelWithData(binding.IntToString(t.gold)),
			newImageFromResource(resourceGoldPng, iconMinSize),
		),
	)
}

// SetProfile sets all labels and images according to the provided Profile
func (t *profileToolbar) SetProfile(p *structs.Profile) error {
	return errors.Join(
		t.title.Set(p.Title),
		t.gold.Set(p.Gold),
		t.ancientCoins.Set(p.AncientCoins),
		t.doubloons.Set(p.Doubloons),
	)
}

func newProfileToolbar() *profileToolbar {
	t := &profileToolbar{
		title:        binding.NewString(),
		gold:         binding.NewInt(),
		doubloons:    binding.NewInt(),
		ancientCoins: binding.NewInt(),
	}
	_ = t.title.Set("<title unknown>")
	return t
}
