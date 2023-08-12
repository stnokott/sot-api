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

	accordion *widget.Accordion
	campaigns *widget.AccordionItem
	emblems   *widget.AccordionItem

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
					s.accordion,
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

	campaigns := widget.NewAccordionItem(
		"Campaigns", container.New(layout.NewFormLayout()),
	)
	emblems := widget.NewAccordionItem(
		"Emblems", container.New(layout.NewFormLayout()),
	)
	accordion := widget.NewAccordion()

	return &summaryView{
		placeholder:       placeholder,
		name:              name,
		motto:             motto,
		rankName:          rankName,
		progressContainer: progressContainer,
		accordion:         accordion,
		campaigns:         campaigns,
		emblems:           emblems,
	}
}

// SetReputation updates the view with new data
func (s *summaryView) SetReputation(name string, rep *structs.Reputation) {
	s.placeholder.Hide()
	s.name.Text = name
	s.motto.Text = "\"" + rep.Motto + "\""

	s.setRankName(rep.RankName)
	s.setUnlockSummaries(rep.UnlockSummaries())
	s.setCampaigns(rep.Campaigns)
	s.setEmblems(rep.Emblems)
	if len(s.accordion.Items) == 1 {
		s.accordion.Open(0)
	} else {
		s.accordion.CloseAll()
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

func (s *summaryView) setCampaigns(campaigns map[string]structs.Campaign) {
	s.accordion.Remove(s.campaigns)
	if campaigns == nil {
		return
	}
	container := s.campaigns.Detail.(*fyne.Container)
	container.RemoveAll()
	items := make([]fyne.CanvasObject, len(campaigns)*2)
	i := 0
	for _, campaign := range campaigns {
		items[i] = widget.NewLabel(campaign.Title)
		subtitle := widget.NewLabel(campaign.Subtitle)
		subtitle.Wrapping = fyne.TextWrapWord
		items[i+1] = subtitle
		i += 2
	}
	container.Objects = items
	container.Refresh()
	s.accordion.Append(s.campaigns)
}

func (s *summaryView) setEmblems(emblems structs.Emblems) {
	s.accordion.Remove(s.emblems)
	if emblems == nil {
		return
	}
	container := s.emblems.Detail.(*fyne.Container)
	container.RemoveAll()
	items := make([]fyne.CanvasObject, len(emblems)*2)
	i := 0
	for _, emblem := range emblems {
		items[i] = widget.NewLabel(emblem.Title)
		subtitle := widget.NewLabel(emblem.Subtitle)
		subtitle.Wrapping = fyne.TextWrapWord
		items[i+1] = subtitle
		i += 2
	}
	container.Objects = items
	container.Refresh()
	s.accordion.Append(s.emblems)
}
