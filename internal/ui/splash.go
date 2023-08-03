package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

func newSplashWindow(a fyne.App) (w fyne.Window) {
	if drv, ok := a.Driver().(desktop.Driver); ok {
		w = drv.CreateSplashWindow()
	} else {
		w = a.NewWindow("Splash Screen")
	}
	w.SetPadded(false)
	splashImg := canvas.NewImageFromResource(resourceSplashPng)
	splashImg.FillMode = canvas.ImageFillOriginal
	w.SetContent(splashImg)
	return
}
