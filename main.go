// Package main runs the application loop
package main

import (
	"github.com/stnokott/sot-api/internal/ui"
	"go.uber.org/zap"
)

func main() {
	// TODO: fix duplicated logger fields
	logger, _ := zap.NewDevelopment(zap.Fields(zap.String("module", "main")))
	//logger, _ := zap.NewProduction()
	defer logger.Sync()

	appLogger := logger.With(zap.String("module", "client"))
	app := ui.NewApp(appLogger)
	app.Run()
}
