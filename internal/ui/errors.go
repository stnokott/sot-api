package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/backend"
)

type errorOverlay struct {
	apiHealthIcon *canvas.Image
	apiReqIcon    *canvas.Image
	apiRespIcon   *canvas.Image
	iconSuccess   fyne.Resource
	iconErr       fyne.Resource
	iconUnknown   fyne.Resource
	lblErrDetails *widget.Label

	refreshInterval time.Duration

	widget.BaseWidget
}

func newErrorOverlay(refreshInterval time.Duration) *errorOverlay {
	iconSuccess := theme.NewPrimaryThemedResource(theme.ConfirmIcon())
	iconErr := theme.NewErrorThemedResource(theme.CancelIcon())
	iconUnknown := theme.NewDisabledResource(theme.QuestionIcon())

	iconSize := fyne.NewSize(32, 32)

	lblErrDetails := widget.NewLabel("n/a")
	lblErrDetails.Wrapping = fyne.TextWrapWord

	return &errorOverlay{
		apiHealthIcon: newImageFromResource(iconSuccess, iconSize),
		apiReqIcon:    newImageFromResource(iconErr, iconSize),
		apiRespIcon:   newImageFromResource(iconSuccess, iconSize),
		iconSuccess:   iconSuccess,
		iconErr:       iconErr,
		iconUnknown:   iconUnknown,

		refreshInterval: refreshInterval,

		lblErrDetails: lblErrDetails,
	}
}

func (o *errorOverlay) CreateRenderer() fyne.WidgetRenderer {
	o.ExtendBaseWidget(o)

	lblTitle := canvas.NewText("Errors occured", color.White)
	lblTitle.TextStyle.Bold = true
	lblTitle.TextSize = 20
	lblSubtitle := canvas.NewText("Updates paused, checking again in "+refreshInterval.String(), color.Gray{200})
	lblSubtitle.TextSize = 12

	return widget.NewSimpleRenderer(
		container.NewMax(
			canvas.NewRectangle(color.NRGBA{R: 50, A: 220}),
			container.NewCenter(
				container.NewVBox(
					lblTitle,
					lblSubtitle,
					container.NewHBox(o.apiHealthIcon, widget.NewLabel("API health")),
					container.NewHBox(o.apiReqIcon, widget.NewLabel("API request")),
					container.NewHBox(o.apiRespIcon, widget.NewLabel("API response decode")),
					container.NewMax(
						canvas.NewRectangle(theme.ErrorColor()),
						o.lblErrDetails,
					),
				),
			),
		),
	)
}

func (o *errorOverlay) setErr(err error) {
	if err == nil {
		o.Hide()
		return
	}

	var resHealth, resReq, resResp fyne.Resource

	switch err.(type) {
	case backend.ErrAPIUnhealthy:
		resHealth = o.iconErr
		resReq = o.iconUnknown
		resResp = o.iconUnknown
	case backend.ErrAPI:
		resHealth = o.iconSuccess
		resReq = o.iconErr
		resResp = o.iconUnknown
	case backend.ErrAPIRespDecode:
		resHealth = o.iconSuccess
		resReq = o.iconSuccess
		resResp = o.iconErr
	}

	o.apiHealthIcon.Resource = resHealth
	o.apiReqIcon.Resource = resReq
	o.apiRespIcon.Resource = resResp
	o.lblErrDetails.Text = err.Error()

	o.Show()
}
