package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

type reputationView struct {
	reputationCategoryViews  map[string]*reputationCategoryView
	reputationCategoriesGrid *fyne.Container

	widget.BaseWidget
}

func (r *reputationView) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)

	return widget.NewSimpleRenderer(
		container.NewMax(
			container.NewVScroll(
				r.reputationCategoriesGrid,
			),
		),
	)
}

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
	repTypeFlameheart       = "Flameheart"
)

var repTypeColors = map[string]color.Color{
	repTypeAthenasFortune:   color.NRGBA{R: 30, G: 81, B: 98, A: 255},
	repTypeHuntersCall:      color.NRGBA{R: 80, G: 96, B: 99, A: 255},
	repTypeGoldHoarders:     color.NRGBA{R: 171, G: 154, B: 45, A: 255},
	repTypeOrderOfSouls:     color.NRGBA{R: 102, G: 46, B: 55, A: 255},
	repTypeMerchantAlliance: color.NRGBA{R: 19, G: 149, B: 178, A: 255},
	repTypeCreatorCrew:      color.NRGBA{R: 70, G: 70, B: 70, A: 255},
	repTypeBilgeRats:        color.NRGBA{R: 49, G: 9, B: 8, A: 255},
	repTypeTallTales:        color.NRGBA{R: 63, G: 64, B: 21, A: 255},
	repTypeReapersBones:     color.NRGBA{R: 202, G: 61, B: 48, A: 255},
	repTypeFlameheart:       color.NRGBA{R: 53, G: 21, B: 14, A: 255},
}

func newReputationView() *reputationView {
	return &reputationView{
		reputationCategoryViews:  make(map[string]*reputationCategoryView),
		reputationCategoriesGrid: container.NewGridWithColumns(3),
	}
}

func (r *reputationView) SetReputation(data structs.Reputations) {
	// TODO: dont remove all, but instead only apply data
	// maybe need to create dedicated struct?
	r.reputationCategoriesGrid.RemoveAll()
	for repName, repVal := range data {
		categoryView, ok := r.reputationCategoryViews[repName]
		if !ok {
			categoryView = newReputationCategoryView(repName, repTypeColors[repName])
			r.reputationCategoryViews[repName] = categoryView
		}
		for repSumName, repSumVal := range repVal.UnlockSummaries() {
			categoryView.AddReputationSummary(reputationProgressSummary{
				Name:     repSumName,
				Total:    repSumVal.Total,
				Unlocked: repSumVal.Unlocked,
			})
		}
		if repVal.Level != nil {
			categoryView.SetLevel(*repVal.Level, *repVal.Progress)
		} else {
			categoryView.HideLevel()
		}
		r.reputationCategoriesGrid.Add(categoryView)
	}

	// remove reputation types not present in data
	for name := range r.reputationCategoryViews {
		if _, ok := data[name]; !ok {
			delete(r.reputationCategoryViews, name)
		}
	}
}

type reputationCategoryView struct {
	nameLabel        *canvas.Text
	summaryContainer *fyne.Container

	levelContainer *fyne.Container
	levelLabel     *widget.Label
	levelProgress  *widget.ProgressBar

	bgColor color.Color

	widget.BaseWidget
}

func (r *reputationCategoryView) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)

	return widget.NewSimpleRenderer(
		container.NewMax(
			canvas.NewHorizontalGradient(r.bgColor, color.Transparent),
			container.NewVBox(
				r.nameLabel,
				r.levelContainer,
				r.summaryContainer,
			),
		),
	)
}

func newProgressTextFormatter(p *widget.ProgressBar) func() string {
	return func() string {
		return fmt.Sprintf("%.0f/%.0f completed", p.Value, p.Max)
	}
}

func newReputationCategoryView(name string, bg color.Color) *reputationCategoryView {
	lblName := canvas.NewText(name, theme.ForegroundColor())
	lblName.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	lblName.TextSize = 28

	levelContainer := container.New(layout.NewFormLayout())
	levelLabel := widget.NewLabel("Level ?")
	levelContainer.Add(levelLabel)
	levelProgress := widget.NewProgressBar()
	levelProgress.Min, levelProgress.Max = 0, 1
	levelContainer.Add(levelProgress)
	levelContainer.Hide()

	return &reputationCategoryView{
		nameLabel:        lblName,
		levelLabel:       levelLabel,
		levelProgress:    levelProgress,
		levelContainer:   levelContainer,
		summaryContainer: container.New(layout.NewFormLayout()),
		bgColor:          bg,
	}
}

type reputationProgressSummary struct {
	Name     string
	Total    int
	Unlocked int
}

func (r *reputationCategoryView) ClearReputationSummaries() {
	r.summaryContainer.RemoveAll()
}

func (r *reputationCategoryView) AddReputationSummary(s reputationProgressSummary) {
	r.summaryContainer.Add(widget.NewLabel(s.Name))
	progress := widget.NewProgressBar()
	progress.TextFormatter = newProgressTextFormatter(progress)
	progress.Min = 0
	progress.Max = float64(s.Total)
	progress.Value = float64(s.Unlocked)
	r.summaryContainer.Add(progress)
}

func (r *reputationCategoryView) SetLevel(level int, progress float64) {
	r.levelLabel.SetText("Level " + strconv.Itoa(level))
	r.levelProgress.SetValue(progress)
	r.levelContainer.Show()
}

func (r *reputationCategoryView) HideLevel() {
	r.levelContainer.Hide()
}
