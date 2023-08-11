package reputation

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

type summaryView struct {
	name              *canvas.Text
	motto             *canvas.Text
	rankName          *widget.Label
	progressContainer *fyne.Container

	widget.BaseWidget
}

func (s *summaryView) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewVBox(
				s.name,
				s.motto,
				s.rankName,
				canvas.NewLine(theme.ForegroundColor()),
				s.progressContainer,
			),
		),
	)
}

func newSummaryView() *summaryView {
	name := canvas.NewText("n/a", theme.ForegroundColor())
	name.TextStyle.Bold = true
	name.TextSize = theme.TextHeadingSize()

	motto := canvas.NewText("n/a", theme.ForegroundColor())
	motto.TextStyle.Italic = true
	motto.TextSize = theme.CaptionTextSize()

	rankName := widget.NewLabel("n/a")

	progressContainer := container.New(layout.NewFormLayout())

	return &summaryView{
		name:              name,
		motto:             motto,
		rankName:          rankName,
		progressContainer: progressContainer,
	}
}

// SetReputation updates the view with new data
func (s *summaryView) SetReputation(name string, rep *structs.Reputation) {
	s.name.Text = name
	s.motto.Text = "\"" + rep.Motto + "\""

	s.setRankName(rep.RankName)
	s.setUnlockSummaries(rep.UnlockSummaries())

	s.Refresh()
}

func (s *summaryView) setRankName(n *string) {
	if n != nil {
		s.rankName.SetText("Rank: " + *n)
		s.rankName.Show()
	} else {
		s.rankName.Hide()
	}
}

func (s *summaryView) setUnlockSummaries(data map[string]structs.UnlockSummary) {
	sumItems := make([]fyne.CanvasObject, len(data)*2)
	i := 0
	for sumName, sumVal := range data {
		sumItems[i] = widget.NewLabel(sumName)
		pb := widget.NewProgressBar()
		pb.Max = float64(sumVal.Total)
		pb.Value = float64(sumVal.Unlocked)
		pb.TextFormatter = newSummaryTextFormatter(pb)
		sumItems[i+1] = pb
		i += 2
	}
	s.progressContainer.Objects = sumItems
}

func newSummaryTextFormatter(p *widget.ProgressBar) func() string {
	return func() string {
		return fmt.Sprintf("%.0f/%.0f completed", p.Value, p.Max)
	}
}
