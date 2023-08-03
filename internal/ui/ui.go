// Package ui provides components for the GUI
package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api"
	"github.com/stnokott/sot-api/internal/backend"
	"go.uber.org/zap"
)

const (
	appTitle string = "Sea of Thieves Tracker"
)

var (
	defaultSize     = fyne.NewSize(800, 600)
	refreshInterval = 30 * time.Second
)

// App coordinates API and UI
type App struct {
	client    *api.Client
	scheduler *backend.Scheduler

	profileToolbar *profileToolbar
	statusBar      *statusBar
	errorOverlay   *errorOverlay

	logger *zap.Logger

	app          fyne.App
	rootWindow   fyne.Window
	splashWindow fyne.Window
}

// NewApp creates a new root app
func NewApp(c *api.Client, logger *zap.Logger) *App {
	a := app.New()
	w := a.NewWindow(appTitle)
	w.SetMaster()

	profile := newProfileToolbar()
	statusBar := newStatusBar()
	errorOverlay := newErrorOverlay(refreshInterval)
	errorOverlay.Hide()

	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		container.NewTabItem("Tab 2", widget.NewLabel("World!")),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(container.NewMax(
		container.NewBorder(
			profile,
			statusBar,
			nil,
			nil,
			tabs,
		),
		errorOverlay,
	))
	w.SetPadded(false)
	w.Resize(defaultSize)

	return &App{
		client:    c,
		scheduler: backend.NewScheduler(c, refreshInterval, logger.With(zap.String("module", "scheduler"))),

		profileToolbar: profile,
		statusBar:      statusBar,
		errorOverlay:   errorOverlay,

		logger: logger,

		app:          a,
		rootWindow:   w,
		splashWindow: newSplashWindow(a),
	}
}

// Run calls ShowAndRun on the underlying root window (blocking)
func (a *App) Run() {
	a.logger.Info("starting app")

	go a.refreshTask()
	a.splashWindow.ShowAndRun()
}

// refreshTask should be ran as goroutine in the background
func (a *App) refreshTask() {
	a.logger.Debug("starting scheduler")
	chStart, chEnd := a.scheduler.Run()

	splashClosed := false

	var onDone func()
	for {
		select {
		case <-chStart:
			onDone = a.statusBar.DoWork()
		case result := <-chEnd:
			if !splashClosed {
				a.rootWindow.Show()
				a.splashWindow.Close()
				splashClosed = true
			}
			a.errorOverlay.setErr(result.Err)
			if result.Err == nil {
				a.profileToolbar.SetProfile(result.Profile)
			}
			if onDone != nil {
				onDone()
			}
		}
	}
}
