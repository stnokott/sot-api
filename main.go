// Package main runs the application loop
package main

import (
	"github.com/stnokott/sot-api/internal/log"
	"github.com/stnokott/sot-api/internal/ui"
)

func main() {
	appLogger := log.ForModule("app")
	app := ui.NewApp(appLogger)
	app.Run()

	if err := log.Sync(); err != nil {
		panic("logger sync failed: " + err.Error())
	}
}
