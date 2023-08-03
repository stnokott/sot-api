package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type statusBar struct {
	lblStatus      *widget.Label
	lblLastUpdated *widget.Label

	widget.BaseWidget
}

func (s *statusBar) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	return widget.NewSimpleRenderer(
		container.NewHBox(
			s.lblStatus,
			layout.NewSpacer(),
			widget.NewLabel("Last updated:"),
			s.lblLastUpdated,
		),
	)
}

const (
	lblStatusIdle     string = "Idle"
	lblStatusUpdating string = "Updating..."
)

func newStatusBar() *statusBar {
	return &statusBar{
		lblStatus:      widget.NewLabel(lblStatusIdle),
		lblLastUpdated: widget.NewLabel("never"),
	}
}

// DoWork indicates that the app is doing something.
// It returns a function that should be called when the work is done.
func (s *statusBar) DoWork() func() {
	s.lblStatus.SetText(lblStatusUpdating)
	return func() {
		s.lblStatus.SetText(lblStatusIdle)
		s.lblLastUpdated.SetText(time.Now().Format(time.TimeOnly))
	}
}
