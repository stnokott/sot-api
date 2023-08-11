// Package ui provides components for the GUI
package ui

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/sot-api/internal/api"
	"github.com/stnokott/sot-api/internal/backend"
	"github.com/stnokott/sot-api/internal/files"
	"github.com/stnokott/sot-api/internal/log"
	"github.com/stnokott/sot-api/internal/ui/reputation"
	"go.uber.org/zap"
	"golang.org/x/text/language"
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
	scheduler *backend.Scheduler

	profileToolbar *profileToolbar
	statusBar      *statusBar
	reputationView *reputation.View
	errorOverlay   *errorOverlay

	logger *zap.Logger

	app          fyne.App
	rootWindow   fyne.Window
	splashWindow fyne.Window
}

// NewApp creates a new root app
func NewApp(logger *zap.Logger) *App {
	a := app.New()
	a.Settings().ThemeVariant()
	w := a.NewWindow(appTitle)
	w.SetMaster()

	profile := newProfileToolbar()
	statusBar := newStatusBar()
	reputationView := reputation.NewView()
	errorOverlay := newErrorOverlay(refreshInterval)
	errorOverlay.Hide()

	tabs := container.NewAppTabs(
		container.NewTabItem("Reputation", reputationView),
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

	token, err := files.ReadToken()
	if err != nil {
		errorOverlay.SetErr(backend.ErrUnauthorized{Err: err})
	}
	client := api.NewClient(token, language.English, log.ForModule("client"))
	scheduler := backend.NewScheduler(client, refreshInterval, log.ForModule("scheduler"))

	return &App{
		scheduler: scheduler,

		profileToolbar: profile,
		statusBar:      statusBar,
		reputationView: reputationView,
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
	a.app.Run()
}

// refreshTask should be ran as goroutine in the background
func (a *App) refreshTask() {
	a.logger.Debug("starting scheduler")
	chTaskBegin, chTaskEnd, chReset := a.scheduler.Run()
	a.errorOverlay.SetFnAuthenticate(func() {
		go a.requestNewToken(chReset)
	})

	a.splashWindow.Show()
	time.Sleep(3 * time.Second) // give splash screen some time to show
	closeSplash := sync.OnceFunc(func() {
		a.rootWindow.Show()
		a.splashWindow.Close()
	})

	var onDone func()
	for {
		select {
		case <-chTaskBegin:
			onDone = a.statusBar.DoWork()
		case result := <-chTaskEnd:
			closeSplash()
			if result.Err == nil {
				a.errorOverlay.Hide()
				a.profileToolbar.SetProfile(result.Profile)
				a.reputationView.SetReputations(result.Reputations)
			} else {
				a.errorOverlay.SetErr(result.Err)
				a.errorOverlay.Show()
			}
			if onDone != nil {
				onDone()
			}
		}
	}
}

func (a *App) requestNewToken(chReset chan<- backend.SchedulerReset) {
	popup := widget.NewModalPopUp(
		widget.NewLabel("Opening browser for login..."),
		a.rootWindow.Canvas(),
	)
	popup.Show()

	resp := <-api.GetAuthFromBrowser()
	if resp.Err != nil {
		a.errorOverlay.SetErr(resp.Err)
	} else {
		chReset <- backend.SchedulerReset{Token: resp.Token}
	}
	popup.Hide()
}
