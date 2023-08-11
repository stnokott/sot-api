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
	apiHealthIcon       *canvas.Image
	apiUnauthorizedIcon *canvas.Image
	apiReqIcon          *canvas.Image
	apiRespIcon         *canvas.Image
	iconSuccess         fyne.Resource
	iconErr             fyne.Resource
	iconUnknown         fyne.Resource
	lblErrDetails       *widget.Label

	btnAuthenticate *widget.Button

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

	btnAuthenticate := widget.NewButton("Authenticate", nil)
	btnAuthenticate.Importance = widget.MediumImportance

	return &errorOverlay{
		apiHealthIcon:       newImageFromResource(iconSuccess, iconSize),
		apiUnauthorizedIcon: newImageFromResource(iconSuccess, iconSize),
		apiReqIcon:          newImageFromResource(iconSuccess, iconSize),
		apiRespIcon:         newImageFromResource(iconSuccess, iconSize),
		iconSuccess:         iconSuccess,
		iconErr:             iconErr,
		iconUnknown:         iconUnknown,
		btnAuthenticate:     btnAuthenticate,

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
					container.NewHBox(o.apiUnauthorizedIcon, widget.NewLabel("API authorization")),
					container.NewHBox(o.apiReqIcon, widget.NewLabel("API request")),
					container.NewHBox(o.apiRespIcon, widget.NewLabel("API response decode")),
					container.NewMax(
						canvas.NewRectangle(theme.ErrorColor()),
						o.lblErrDetails,
					),
					o.btnAuthenticate,
				),
			),
		),
	)
}

func (o *errorOverlay) SetFnAuthenticate(f func()) {
	o.btnAuthenticate.OnTapped = f
}

func (o *errorOverlay) SetErr(err error) {
	if err == nil {
		return
	}

	o.btnAuthenticate.Hide()
	var resHealth, resAuth, resReq, resResp fyne.Resource

	switch err.(type) {
	case backend.ErrAPIUnhealthy:
		resHealth = o.iconErr
		resAuth = o.iconUnknown
		resReq = o.iconUnknown
		resResp = o.iconUnknown
	case backend.ErrUnauthorized:
		resHealth = o.iconSuccess
		resAuth = o.iconErr
		resReq = o.iconUnknown
		resResp = o.iconUnknown
		o.btnAuthenticate.Show()
	case backend.ErrAPI:
		resHealth = o.iconSuccess
		resAuth = o.iconSuccess
		resReq = o.iconErr
		resResp = o.iconUnknown
	case backend.ErrAPIRespDecode:
		resHealth = o.iconSuccess
		resAuth = o.iconSuccess
		resReq = o.iconSuccess
		resResp = o.iconErr
	default:
		resHealth = o.iconUnknown
		resAuth = o.iconUnknown
		resReq = o.iconUnknown
		resResp = o.iconUnknown
	}

	o.apiHealthIcon.Resource = resHealth
	o.apiUnauthorizedIcon.Resource = resAuth
	o.apiReqIcon.Resource = resReq
	o.apiRespIcon.Resource = resResp
	o.lblErrDetails.Text = err.Error()
}
