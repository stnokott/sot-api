package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
)

func newImageFromResource(r fyne.Resource, minSize fyne.Size) *canvas.Image {
	img := canvas.NewImageFromResource(r)
	img.SetMinSize(minSize)
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScaleFastest
	return img
}

func newImageFromURL(url string, minSize fyne.Size) (*canvas.Image, error) {
	uri, err := storage.ParseURI(url)
	if err != nil {
		return nil, err
	}
	img := canvas.NewImageFromURI(uri)
	img.SetMinSize(minSize)
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScaleFastest
	return img, nil
}
