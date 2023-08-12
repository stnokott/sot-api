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
	// placeholder will be visible until a reputation is set
	placeholder *fyne.Container

	name              *canvas.Text
	motto             *canvas.Text
	rankName          *widget.Label
	progressContainer *fyne.Container

	detailLabel     *canvas.Text
	detailContainer *fyne.Container

	widget.BaseWidget
}

func (s *summaryView) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewVScroll(
				container.NewVBox(
					s.name,
					s.motto,
					s.rankName,
					canvas.NewLine(theme.ForegroundColor()),
					s.progressContainer,
					canvas.NewLine(theme.ForegroundColor()),
					s.detailLabel,
					s.detailContainer,
				),
			),
			s.placeholder,
		),
	)
}

func newSummaryView() *summaryView {
	placeholder := container.NewMax(
		canvas.NewRectangle(theme.BackgroundColor()),
		container.NewCenter(
			widget.NewLabel("Select a reputation on the left"),
		),
	)

	name := canvas.NewText("n/a", theme.ForegroundColor())
	name.TextStyle.Bold = true
	name.TextSize = theme.TextHeadingSize()
	motto := canvas.NewText("n/a", theme.ForegroundColor())
	motto.TextStyle.Italic = true
	motto.TextSize = theme.CaptionTextSize()
	rankName := widget.NewLabel("n/a")

	progressContainer := container.New(layout.NewFormLayout())

	detailLabel := canvas.NewText("n/a", theme.ForegroundColor())
	detailLabel.TextSize = theme.TextSubHeadingSize()
	detailLabel.TextStyle.Bold = true
	detailContainer := container.New(layout.NewFormLayout())

	return &summaryView{
		placeholder:       placeholder,
		name:              name,
		motto:             motto,
		rankName:          rankName,
		progressContainer: progressContainer,
		detailLabel:       detailLabel,
		detailContainer:   detailContainer,
	}
}

// SetReputation updates the view with new data
func (s *summaryView) SetReputation(name string, rep *structs.Reputation) {
	s.placeholder.Hide()
	s.name.Text = name
	s.motto.Text = "\"" + rep.Motto + "\""

	s.setRankName(rep.RankName)
	s.setUnlockSummaries(rep.UnlockSummaries())
	switch {
	case rep.Campaigns != nil:
		s.setCampaigns(rep.Campaigns)
	case rep.Emblems != nil:
		s.setEmblems(rep.Emblems)
	default:
		s.detailLabel.Hide()
		s.detailContainer.Hide()
	}
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

// setCampaigns updates the view with campaign details.
// Should not be called together with setEmblems as only of the two is possible for the same dataset.
func (s *summaryView) setCampaigns(campaigns map[string]structs.Campaign) {
	s.detailLabel.Text = "Campaigns"
	s.detailContainer.RemoveAll()
	items := make([]fyne.CanvasObject, len(campaigns)*2)
	i := 0
	for _, campaign := range campaigns {
		items[i] = widget.NewLabel(campaign.Title)
		subtitle := widget.NewLabel(campaign.Subtitle)
		subtitle.Wrapping = fyne.TextWrapWord
		items[i+1] = subtitle
		i += 2
	}
	s.detailContainer.Objects = items
	s.detailContainer.Refresh()

	s.detailLabel.Show()
	s.detailContainer.Show()
}

// setEmblems updates the view with emblem details.
// Should not be called together with setCampaigns as only of the two is possible for the same dataset.
func (s *summaryView) setEmblems(emblems structs.Emblems) {
	s.detailLabel.Text = "Emblems"
	s.detailContainer.RemoveAll()
	items := make([]fyne.CanvasObject, len(emblems)*2)
	for i, emblem := range emblems {
		items[i*2] = widget.NewLabel(emblem.Title)
		subtitle := widget.NewLabel(emblem.Subtitle)
		subtitle.Wrapping = fyne.TextWrapWord
		items[i*2+1] = subtitle
	}
	s.detailContainer.Objects = items
	s.detailContainer.Refresh()

	s.detailLabel.Show()
	s.detailContainer.Show()
}
