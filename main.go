// Package main runs the application loop
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stnokott/sot-api/internal/api"
	"github.com/stnokott/sot-api/internal/api/structs"
	"github.com/stnokott/sot-api/internal/io"
	"golang.org/x/exp/slog"
	"golang.org/x/text/language"
)

func main() {
	logLevel := slog.LevelDebug
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	rootLogger := slog.New(logHandler).With("module", "main")

	token, err := io.ReadToken()
	if err != nil {
		rootLogger.Error(err.Error())
		panic(err)
	}
	rootLogger.Info("token read from file")

	clientLogger := slog.New(logHandler).With("module", "api-client")
	client := api.NewClient(token, language.German, clientLogger)
	health, err := client.GetHealth()
	if err != nil {
		panic(err)
	}
	if health.HasFailures {
		fmt.Println(health.CurrentStatus)
		panic("SoT API not healthy")
	}

	season, err := client.GetSeason()
	if err != nil {
		panic(err)
	}
	printSeasonChallenges(&season)
}

func printSeasonChallenges(s *structs.Season) {
	for _, group := range s.ChallengeGroups {
		log.Println()
		log.Printf("%s (%d/%d)", group.Title, group.ChallengesCompleted, len(group.Challenges))
		if group.IsCompleted {
			continue
		}
		for _, challenge := range group.Challenges {
			if !challenge.IsCompleted {
				log.Printf("  - %s (Size %v, %d/%d)\n", challenge.Title, challenge.XPGain, challenge.GoalsCompleted, len(challenge.Goals))
				for _, goal := range challenge.Goals {
					if goal.ProgressCurrent != goal.ProgressMax {
						log.Printf("      > %s (Size %v, %d/%d)\n", goal.Title, goal.XPGain, goal.ProgressCurrent, goal.ProgressMax)
					}
				}
			}
		}
	}
}
