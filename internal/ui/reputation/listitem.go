package reputation

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

type coloredListItem struct {
	background *canvas.LinearGradient
	name       *widget.Label

	level *widget.Label

	widget.BaseWidget
}

func newColoredListItem() *coloredListItem {
	background := canvas.NewHorizontalGradient(
		color.Transparent,
		color.Transparent,
	)
	name := widget.NewLabel("n/a")
	level := widget.NewLabel("n/a")
	r := &coloredListItem{
		background: background,
		name:       name,
		level:      level,
	}
	r.ExtendBaseWidget(r)
	return r
}

func (li *coloredListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewMax(
			li.background,
			container.NewHBox(
				li.name,
				layout.NewSpacer(),
				li.level,
			),
		),
	)
}

func (li *coloredListItem) SetReputation(name string, color color.Color, rep *structs.Reputation) {
	li.background.StartColor = color
	li.name.SetText(name)
	if rep.Level != nil {
		li.level.SetText(fmt.Sprintf("Level %d (%.0f%%)", *rep.Level, *rep.Progress*100))
		li.level.Show()
	} else {
		li.level.Hide()
	}
}
