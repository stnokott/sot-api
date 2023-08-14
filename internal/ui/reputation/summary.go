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
	rankName          *canvas.Text
	progressContainer *fyne.Container

	detailContainer *widget.Accordion

	widget.BaseWidget
}

func (s *summaryView) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)

	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewVScroll(
				container.NewVBox(
					container.NewHBox(
						container.NewVBox(
							s.name,
							s.motto,
						),
						layout.NewSpacer(),
						container.NewPadded(s.rankName),
					),
					canvas.NewLine(theme.ForegroundColor()),
					s.progressContainer,
					canvas.NewLine(theme.ForegroundColor()),
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
	rankName := canvas.NewText("n/a", theme.ForegroundColor())
	rankName.TextSize = theme.TextSubHeadingSize()

	progressContainer := container.New(layout.NewFormLayout())

	detailContainer := widget.NewAccordion()

	return &summaryView{
		placeholder:       placeholder,
		name:              name,
		motto:             motto,
		rankName:          rankName,
		progressContainer: progressContainer,
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
		s.detailContainer.Hide()
	}
	s.Refresh()
}

func (s *summaryView) setRankName(n *string) {
	if n != nil {
		s.rankName.Text = "Rank: " + *n
		s.rankName.Show()
	} else {
		s.rankName.Hide()
	}
}

func (s *summaryView) setUnlockSummaries(data map[string]structs.UnlockSummary) {
	var sumItems []fyne.CanvasObject
	i := 0
	for sumName, sumVal := range data {
		if sumVal.Total == 0 {
			continue
		}
		sumItems = append(sumItems, widget.NewLabel(sumName))
		pb := widget.NewProgressBar()
		pb.Max = float64(sumVal.Total)
		pb.Value = float64(sumVal.Unlocked)
		pb.TextFormatter = newSummaryTextFormatter(pb)
		sumItems = append(sumItems, pb)
		i += 2
	}
	s.progressContainer.Objects = sumItems
}

func newSummaryTextFormatter(p *widget.ProgressBar) func() string {
	return func() string {
		return fmt.Sprintf("%.0f/%.0f completed", p.Value, p.Max)
	}
}

func (s *summaryView) setEmblemMap(emblemMap map[string][]structs.Emblem) {
	items := make([]*widget.AccordionItem, len(emblemMap))
	i := 0
	for name, emblems := range emblemMap {
		cont := container.New(layout.NewFormLayout())
		contItems := make([]fyne.CanvasObject, 2*len(emblems))
		for j, emblem := range emblems {
			contItems[j*2] = widget.NewLabel(emblem.Title)
			subtitle := widget.NewLabel(emblem.Subtitle)
			subtitle.Wrapping = fyne.TextWrapWord
			contItems[j*2+1] = subtitle
		}
		cont.Objects = contItems
		// TODO: needed?
		// cont.Refresh()
		items[i] = widget.NewAccordionItem(name, cont)
		i++
	}
	s.detailContainer.Items = items
	if len(s.detailContainer.Items) == 1 {
		s.detailContainer.Open(0)
	}
	s.detailContainer.Show()
	s.detailContainer.Refresh()
}

// setCampaigns updates the view with campaign details.
// Should not be called together with setEmblems as only of the two is possible for the same dataset.
func (s *summaryView) setCampaigns(campaigns map[string]structs.Campaign) {
	emblemMap := make(map[string][]structs.Emblem, len(campaigns))
	for _, campaign := range campaigns {
		emblemMap[campaign.Title] = campaign.Emblems
	}
	s.setEmblemMap(emblemMap)
}

// setEmblems updates the view with emblem details.
// Should not be called together with setCampaigns as only of the two is possible for the same dataset.
func (s *summaryView) setEmblems(emblems structs.Emblems) {
	emblemMap := map[string][]structs.Emblem{
		"Emblems": emblems,
	}
	s.setEmblemMap(emblemMap)
}
