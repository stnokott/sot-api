// Package backend schedules API calls and handles results
package backend

import (
	"errors"
	"reflect"
	"time"

	"github.com/stnokott/sot-api/internal/api"
	"github.com/stnokott/sot-api/internal/api/structs"
	"go.uber.org/zap"
)

// Scheduler handles API calls and bundles results and errors.
type Scheduler struct {
	client *api.Client

	refreshInterval time.Duration

	logger *zap.Logger
}

// NewScheduler returns a new instance of Scheduler.
func NewScheduler(client *api.Client, refreshInterval time.Duration, logger *zap.Logger) *Scheduler {
	return &Scheduler{
		client:          client,
		refreshInterval: refreshInterval,
		logger:          logger,
	}
}

// JobResult contains the result of a scheduled job.
// If Err is not nil, all other fields will be nil
type JobResult struct {
	Profile     *structs.Profile
	Reputations structs.Reputations

	Err error
}

// ErrAPIUnhealthy occurs when the API is unhealthy
type ErrAPIUnhealthy error

// ErrAPI occurs when there is an error querying the API
type ErrAPI error

// ErrAPIRespDecode occurs when an API response is received, but it could not
// be decoded.
type ErrAPIRespDecode error

// Run starts the scheduler and returns one channel for beginning a task and one for finishing it.
// It will run forever, channels will never be closed.
func (s *Scheduler) Run() (start <-chan struct{}, end <-chan JobResult) {
	chStart := make(chan struct{})
	chEnd := make(chan JobResult)
	start, end = chStart, chEnd

	go func() {
		s.logger.Debug("running initial update task")
		chStart <- struct{}{}
		chEnd <- s.doTask()
		s.logger.Debug("finished initial task")
		// goroutine will never end, so no need to use NewTicker (ok to "leak")
		for range time.Tick(s.refreshInterval) {
			s.logger.Debug("starting update task")
			chStart <- struct{}{}
			chEnd <- s.doTask()
			s.logger.Debug("finished update task")
		}
	}()

	s.logger.Debug("scheduler started")
	return
}

func (s *Scheduler) doTask() (r JobResult) {
	var err error
	defer func() {
		r.Err = convertAPIErr(err)
		if r.Err != nil {
			s.logger.Debug("got " + reflect.TypeOf(r.Err).String() + " error, checking API health")
			if health, err := s.client.GetHealth(); err != nil {
				s.logger.Debug("could not retrieve API health, falling back to original error")
				r.Err = errors.Join(r.Err, err)
			} else if health.HasFailures {
				s.logger.Debug("API is unhealthy, overwriting error")
				r.Err = errors.New(health.String()).(ErrAPIUnhealthy)
			}
		}
	}()

	s.logger.Debug("getting profile")
	r.Profile, err = s.client.GetProfile()
	if err != nil {
		return
	}
	r.Reputations, err = s.client.GetReputation()

	return
}

// convertAPIErr converts an error returned by the API to an error defined in this package.
func convertAPIErr(err error) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case api.ErrHTTP:
		return err.(ErrAPI)
	case api.ErrResponseDecode:
		return err.(ErrAPIRespDecode)
	default:
		return err
	}
}
