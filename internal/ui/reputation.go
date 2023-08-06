package ui

import (
	"fmt"
	"image/color"

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
	// TODO: remove reputation types not present in data
	r.reputationCategoriesGrid.RemoveAll()
	for repName, repVal := range data {
		categoryView, ok := r.reputationCategoryViews[repName]
		if !ok {
			categoryView = newReputationCategoryView(repName, repTypeColors[repName])
		}
		if repVal.EmblemsTotal != nil {
			categoryView.AddReputationSummary(reputationProgressSummary{
				Name:     "Emblems",
				Total:    *repVal.EmblemsTotal,
				Unlocked: *repVal.EmblemsUnlocked,
			})
		}
		if repVal.ItemsTotal != nil {
			categoryView.AddReputationSummary(reputationProgressSummary{
				Name:     "Items",
				Total:    *repVal.ItemsTotal,
				Unlocked: *repVal.ItemsUnlocked,
			})
		}
		if repVal.TitlesTotal != nil {
			categoryView.AddReputationSummary(reputationProgressSummary{
				Name:     "Titles",
				Total:    *repVal.TitlesTotal,
				Unlocked: *repVal.TitlesUnlocked,
			})
		}
		if repVal.PromotionsTotal != nil {
			categoryView.AddReputationSummary(reputationProgressSummary{
				Name:     "Promotions",
				Total:    *repVal.PromotionsTotal,
				Unlocked: *repVal.PromotionsUnlocked,
			})
		}
		// TODO: split in separate function
		r.reputationCategoriesGrid.Add(categoryView)
	}
}

type reputationCategoryView struct {
	lblName     *canvas.Text
	summaryForm *fyne.Container
	bgColor     color.Color

	widget.BaseWidget
}

func (r *reputationCategoryView) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)

	// TODO: display level

	return widget.NewSimpleRenderer(
		container.NewMax(
			canvas.NewHorizontalGradient(r.bgColor, color.Transparent),
			container.NewVBox(
				r.lblName,
				r.summaryForm,
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

	return &reputationCategoryView{
		lblName:     lblName,
		summaryForm: container.New(layout.NewFormLayout()),
		bgColor:     bg,
	}
}

type reputationProgressSummary struct {
	Name     string
	Total    int
	Unlocked int
}

func (r *reputationCategoryView) ClearReputationSummaries() {
	r.summaryForm.RemoveAll()
}

func (r *reputationCategoryView) AddReputationSummary(s reputationProgressSummary) {
	r.summaryForm.Add(widget.NewLabel(s.Name))
	progress := widget.NewProgressBar()
	progress.TextFormatter = newProgressTextFormatter(progress)
	progress.Min = 0
	progress.Max = float64(s.Total)
	progress.Value = float64(s.Unlocked)
	r.summaryForm.Add(progress)
}
