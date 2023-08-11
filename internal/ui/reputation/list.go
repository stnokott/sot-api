// Package reputation contains UI objects for displaying reputation in various categories
package reputation

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api/structs"
)

// View is a list of reputations together with a detailed view of the selected reputation
type View struct {
	reputationNames        []string
	reputationNamesBinding binding.ExternalStringList
	reputationItems        structs.Reputations

	summaryView *summaryView

	widget.BaseWidget
}

// NewView creates a new ReputationView instance
func NewView() *View {
	r := &View{
		summaryView: newSummaryView(),
	}
	r.reputationNamesBinding = binding.BindStringList(&r.reputationNames)
	return r
}

// CreateRenderer implements the fyne.Widget interface
func (v *View) CreateRenderer() fyne.WidgetRenderer {
	v.ExtendBaseWidget(v)

	reputationList := widget.NewListWithData(
		v.reputationNamesBinding,
		func() fyne.CanvasObject {
			return newColoredListItem()
		},
		v.updateReputationItem,
	)
	reputationList.OnSelected = v.onItemSelected

	return widget.NewSimpleRenderer(
		container.NewHSplit(
			reputationList,
			v.summaryView,
		),
	)
}

func (v *View) updateReputationItem(di binding.DataItem, co fyne.CanvasObject) {
	name, err := di.(binding.String).Get()
	if err != nil {
		// TODO: display error
		panic(err)
	}
	displayName := repTypeDisplaynames[name]
	color := repTypeColors[name]
	co.(*coloredListItem).SetReputation(displayName, color, v.reputationItems[name])
}

// SetReputations can be used to update the reputations
func (v *View) SetReputations(data structs.Reputations) {
	// get reputation names
	names := make([]string, len(data))
	i := 0
	for repName := range data {
		names[i] = repName
		i++
	}
	slices.Sort(names) // sort names

	v.reputationNames = names
	if err := v.reputationNamesBinding.Reload(); err != nil {
		// should not happen
		panic(err)
	}
	// TODO: update summary for currently selected item
	v.reputationItems = data
}

func (v *View) onItemSelected(i int) {
	itemName := v.reputationNames[i]
	displayName := repTypeDisplaynames[itemName]
	v.summaryView.SetReputation(displayName, v.reputationItems[itemName])
}
