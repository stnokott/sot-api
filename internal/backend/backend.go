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

// ErrAPIUnhealthy occurs when the API is unhealthy
type ErrAPIUnhealthy struct{}

func (e ErrAPIUnhealthy) Error() string {
	return "API is unhealthy"
}

// ErrUnauthorized occurs when authorization using the provided RAT token did not succeed
type ErrUnauthorized struct {
	Err error
}

func (e ErrUnauthorized) Error() string {
	return "RAT token is expired or invalid"
}

func (e ErrUnauthorized) Unwrap() error {
	return e.Err
}

// ErrAPI occurs when there is an error querying the API
type ErrAPI struct {
	APIErr error
}

func (e ErrAPI) Error() string {
	return "there was an error querying the API: " + e.APIErr.Error()
}

func (e ErrAPI) Unwrap() error {
	return e.APIErr
}

// ErrAPIRespDecode occurs when an API response is received, but it could not
// be decoded.
type ErrAPIRespDecode struct {
	DecodeErr error
}

func (e ErrAPIRespDecode) Error() string {
	return "there was an error decoding the API response: " + e.DecodeErr.Error()
}

func (e ErrAPIRespDecode) Unwrap() error {
	return e.DecodeErr
}

// JobResult contains the result of a scheduled job.
// If Err is not nil, all other fields will be nil
type JobResult struct {
	Profile     *structs.Profile
	Reputations structs.Reputations

	Err error
}

// SchedulerReset can be sent to a running scheduler via channel to update the RAT token while running
type SchedulerReset struct {
	Token string
}

// Run starts the scheduler and returns one channel for beginning a task and one for finishing it.
// It will run forever, channels will never be closed.
func (s *Scheduler) Run() (start <-chan struct{}, end <-chan JobResult, reset chan<- SchedulerReset) {
	chStart := make(chan struct{})
	chEnd := make(chan JobResult)
	chReset := make(chan SchedulerReset)
	start, end, reset = chStart, chEnd, chReset

	doTask := func() {
		s.logger.Debug("running task")
		chStart <- struct{}{}
		chEnd <- s.getAPIData()
		s.logger.Debug("finished task")
	}

	go func() {
		doTask() // initial run
		ticker := time.NewTicker(s.refreshInterval)
		for {
			select {
			case <-ticker.C:
				doTask()
			case reset := <-chReset:
				s.logger.Debug("reset requested")
				if reset.Token != "" {
					s.client.SetToken(reset.Token)
					s.logger.Debug("token updated")
				}
				doTask()
				ticker.Reset(s.refreshInterval)
			}
		}
	}()

	s.logger.Debug("scheduler started")
	return
}

func (s *Scheduler) getAPIData() (r JobResult) {
	var err error
	defer func() {
		r.Err = convertAPIErr(err)
		if r.Err != nil {
			s.logger.Debug("got " + reflect.TypeOf(r.Err).String() + " error, checking API health")
			if health, err := s.client.GetHealth(); err != nil {
				s.logger.Debug("could not retrieve API health, falling back to original error")
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
	switch err := err.(type) {
	case api.ErrHTTP:
		return ErrAPI{APIErr: errors.Unwrap(err)}
	case api.ErrUnauthorized:
		return ErrUnauthorized{Err: errors.Unwrap(err)}
	case api.ErrResponseDecode:
		return ErrAPIRespDecode{DecodeErr: errors.Unwrap(err)}
	default:
		return err
	}
}
