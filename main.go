// Package main runs the application loop
package main

import (
	"github.com/stnokott/sot-api/internal/api"
	"github.com/stnokott/sot-api/internal/io"
	"github.com/stnokott/sot-api/internal/ui"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func main() {
	logger, _ := zap.NewDevelopment(zap.Fields(zap.String("module", "main")))
	//logger, _ := zap.NewProduction()
	defer logger.Sync()

	token, err := io.ReadToken()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	logger.Info("token read from file")

	clientLogger := logger.With(zap.String("module", "client"))
	client := api.NewClient(token, language.German, clientLogger)

	appLogger := logger.With(zap.String("module", "client"))
	app := ui.NewApp(client, appLogger)
	app.Run()
}
